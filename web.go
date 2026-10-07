package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"pmscanner/fx"
	"pmscanner/store"
)

// server rescans in the background and serves scans from the database.
type server struct {
	db       *store.Store
	mu       sync.RWMutex
	rates    fx.Rates
	ratesAt  time.Time
	scanning bool
	lastRun  time.Time // last scan attempt, even if it could not be saved
}

func runServer(addr string, refresh time.Duration, db *store.Store) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := &server{db: db}
	if refresh > 0 {
		go s.scanEvery(refresh)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handlePage)
	mux.HandleFunc("GET /api/prices", s.handlePrices)
	mux.HandleFunc("GET /api/scans", s.handleScans)
	url := "http://" + displayAddr(ln.Addr())
	if refresh > 0 {
		log.Printf("serving on %s (rescan every %s)", url, refresh)
	} else {
		log.Printf("serving on %s (automatic rescans off)", url)
	}
	return http.Serve(ln, mux)
}

// displayAddr shows a listen address as a URL host, using localhost when
// listening on all interfaces.
func displayAddr(addr net.Addr) string {
	tcp, ok := addr.(*net.TCPAddr)
	if !ok || !tcp.IP.IsUnspecified() {
		return addr.String()
	}
	return net.JoinHostPort("localhost", strconv.Itoa(tcp.Port))
}

// scanEvery rescans at the given interval, counting from the latest saved
// scan so that restarting the server does not trigger a needless scan.
func (s *server) scanEvery(interval time.Duration) {
	for {
		if wait := interval - s.sinceLastScan(); wait > 0 {
			log.Printf("next scan in %s", wait.Round(time.Minute))
			time.Sleep(wait)
			continue
		}
		s.rescan()
	}
}

// sinceLastScan returns how long ago the latest scan ran, or a very long
// time when there is none.
func (s *server) sinceLastScan() time.Duration {
	s.mu.RLock()
	last := s.lastRun
	s.mu.RUnlock()
	scans, err := s.db.Scans(1)
	if err != nil {
		log.Printf("read last scan: %v", err)
	}
	if len(scans) > 0 && scans[0].At.After(last) {
		last = scans[0].At
	}
	if last.IsZero() {
		return 1<<63 - 1
	}
	return time.Since(last)
}

func (s *server) rescan() {
	s.setScanning(true)
	defer s.setScanning(false)

	start := time.Now()
	s.mu.Lock()
	s.lastRun = start
	s.mu.Unlock()
	entries := scan()
	log.Printf("scanned %d products in %s, %d failed", len(entries), time.Since(start).Round(time.Second), failures(entries))
	if _, err := s.db.Save(start, toPrices(entries)); err != nil {
		log.Printf("save scan: %v", err)
	}
}

// ratesMaxAge is how long exchange rates are reused; they change once a day.
const ratesMaxAge = 6 * time.Hour

// currentRates returns cached exchange rates, refreshing them when stale.
func (s *server) currentRates() fx.Rates {
	s.mu.RLock()
	rates, stale := s.rates, time.Since(s.ratesAt) > ratesMaxAge
	s.mu.RUnlock()
	if !stale {
		return rates
	}
	// Fetched outside the lock; concurrent requests may both fetch, which is harmless.
	fresh, err := fx.FetchRates()
	if err != nil {
		log.Printf("exchange rates: %v", err)
		return rates
	}
	s.mu.Lock()
	s.rates, s.ratesAt = fresh, time.Now()
	s.mu.Unlock()
	return fresh
}

func (s *server) setScanning(v bool) {
	s.mu.Lock()
	s.scanning = v
	s.mu.Unlock()
}

func (s *server) isScanning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scanning
}

// serverLinks points to the live server's routes.
type serverLinks struct{}

func (serverLinks) Scan(id int64, latest bool) string {
	if latest {
		return "/"
	}
	return "/?scan=" + strconv.FormatInt(id, 10)
}
func (serverLinks) JSON(id int64) string { return "/api/prices?scan=" + strconv.FormatInt(id, 10) }
func (serverLinks) Latest() string       { return "/" }

// scanID reads ?scan=ID, defaulting to the latest scan.
func (s *server) scanID(r *http.Request) (int64, error) {
	if v := r.URL.Query().Get("scan"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, store.ErrNotFound
		}
		return id, nil
	}
	return s.db.Latest()
}

func (s *server) handlePage(w http.ResponseWriter, r *http.Request) {
	id, err := s.scanID(r)
	var data pageData
	if err == nil {
		data, err = buildPage(s.db, id, serverLinks{}, s.currentRates())
	}
	switch {
	case errors.Is(err, store.ErrNotFound) && !r.URL.Query().Has("scan"):
		data = pageData{} // no scan yet
	case errors.Is(err, store.ErrNotFound):
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	case err != nil:
		s.serverError(w, err)
		return
	}
	data.Scanning = s.isScanning()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := renderPage(w, data); err != nil {
		log.Printf("render: %v", err)
	}
}

func (s *server) handlePrices(w http.ResponseWriter, r *http.Request) {
	id, err := s.scanID(r)
	if err != nil {
		s.apiError(w, err)
		return
	}
	v, err := pricesJSON(s.db, id)
	if err != nil {
		s.apiError(w, err)
		return
	}
	writeJSON(w, v)
}

func (s *server) handleScans(w http.ResponseWriter, r *http.Request) {
	scans, err := s.db.Scans(scanListLimit)
	if err != nil {
		s.serverError(w, err)
		return
	}
	writeJSON(w, scansJSON(scans))
}

func (s *server) serverError(w http.ResponseWriter, err error) {
	log.Printf("request: %v", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func (s *server) apiError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	}
	s.serverError(w, err)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}
