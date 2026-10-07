package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"pmscanner/store"
)

//go:embed web.html
var pageHTML string

var page = template.Must(template.New("page").Parse(pageHTML))

// categories group products on the page, matched by description prefix.
var categories = []string{
	"Gold bar 1kg",
	"Gold Krugerrand 1oz",
	"Gold Maple Leaf 1oz",
	"Silver Maple Leaf 1oz",
}

const otherCategory = "Other"

// scanListLimit caps how many past scans the page and API list.
const scanListLimit = 500

func categoryOf(description string) string {
	for _, c := range categories {
		if strings.HasPrefix(description, c) {
			return c
		}
	}
	return otherCategory
}

// server rescans in the background and serves scans from the database.
type server struct {
	db       *store.Store
	mu       sync.RWMutex
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
	if refresh > 0 {
		log.Printf("serving on http://%s (rescan every %s)", ln.Addr(), refresh)
	} else {
		log.Printf("serving on http://%s (automatic rescans off)", ln.Addr())
	}
	return http.Serve(ln, mux)
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

type row struct {
	Site, Description, URL, Price, PerGram, Currency, Error string
	Best                                                    bool
}

type group struct {
	Name string
	Rows []row
}

type scanOption struct {
	ID       int64
	Label    string
	Selected bool
}

type pageData struct {
	Groups   []group
	Failed   []row
	Scanned  string
	Scanning bool
	IsLatest bool
	Scans    []scanOption
}

const timeFormat = "2006-01-02 15:04"

func (s *server) handlePage(w http.ResponseWriter, r *http.Request) {
	data := pageData{Scanning: s.isScanning()}

	id, err := s.scanID(r)
	switch {
	case errors.Is(err, store.ErrNotFound) && r.URL.Query().Has("scan"):
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	case errors.Is(err, store.ErrNotFound):
		s.render(w, data) // no scan yet
		return
	case err != nil:
		s.serverError(w, err)
		return
	}

	scanned, prices, err := s.db.Prices(id)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	}
	if err != nil {
		s.serverError(w, err)
		return
	}
	scans, err := s.db.Scans(scanListLimit)
	if err != nil {
		s.serverError(w, err)
		return
	}

	data.Scanned = scanned.Local().Format(timeFormat)
	data.IsLatest = len(scans) > 0 && scans[0].ID == id
	for _, sc := range scans {
		label := sc.At.Local().Format(timeFormat)
		if sc.Failed > 0 {
			label += " (" + strconv.Itoa(sc.Failed) + " failed)"
		}
		data.Scans = append(data.Scans, scanOption{ID: sc.ID, Label: label, Selected: sc.ID == id})
	}

	byName := map[string]*group{}
	for _, name := range append(categories, otherCategory) {
		data.Groups = append(data.Groups, group{Name: name})
	}
	for i := range data.Groups {
		byName[data.Groups[i].Name] = &data.Groups[i]
	}
	// Prices are saved sorted by price per gram.
	for _, p := range prices {
		rw := row{Site: p.Site, Description: p.Description, URL: p.URL}
		if p.Error != "" {
			rw.Error = p.Error
			data.Failed = append(data.Failed, rw)
			continue
		}
		rw.Price, rw.PerGram, rw.Currency = p.Price.StringFixed(2), p.EURPerGram.StringFixed(2), p.Currency
		g, ok := byName[p.Category]
		if !ok {
			g = byName[otherCategory]
		}
		rw.Best = len(g.Rows) == 0
		g.Rows = append(g.Rows, rw)
	}
	s.render(w, data)
}

func (s *server) render(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Execute(w, data); err != nil {
		log.Printf("render: %v", err)
	}
}

func (s *server) serverError(w http.ResponseWriter, err error) {
	log.Printf("request: %v", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func (s *server) handlePrices(w http.ResponseWriter, r *http.Request) {
	id, err := s.scanID(r)
	if err != nil {
		s.apiError(w, err)
		return
	}
	scanned, prices, err := s.db.Prices(id)
	if err != nil {
		s.apiError(w, err)
		return
	}
	type item struct {
		Site        string  `json:"site"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
		URL         string  `json:"url"`
		Price       float64 `json:"price,omitempty"`
		Currency    string  `json:"currency,omitempty"`
		EURPerGram  float64 `json:"eurPerGram,omitempty"`
		Error       string  `json:"error,omitempty"`
	}
	items := make([]item, 0, len(prices))
	for _, p := range prices {
		it := item{Site: p.Site, Category: p.Category, Description: p.Description, URL: p.URL, Error: p.Error}
		if p.Error == "" {
			it.Price, _ = p.Price.Float64()
			it.EURPerGram, _ = p.EURPerGram.Round(4).Float64()
			it.Currency = p.Currency
		}
		items = append(items, it)
	}
	writeJSON(w, map[string]any{"scan": id, "scanned": scanned, "prices": items})
}

func (s *server) handleScans(w http.ResponseWriter, r *http.Request) {
	scans, err := s.db.Scans(scanListLimit)
	if err != nil {
		s.serverError(w, err)
		return
	}
	type item struct {
		ID     int64     `json:"id"`
		At     time.Time `json:"at"`
		Count  int       `json:"count"`
		Failed int       `json:"failed"`
	}
	items := make([]item, 0, len(scans))
	for _, sc := range scans {
		items = append(items, item{sc.ID, sc.At, sc.Count, sc.Failed})
	}
	writeJSON(w, map[string]any{"scans": items})
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
