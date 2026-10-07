package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
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

func categoryOf(description string) string {
	for _, c := range categories {
		if strings.HasPrefix(description, c) {
			return c
		}
	}
	return otherCategory
}

// server keeps the latest scan and refreshes it in the background.
type server struct {
	mu       sync.RWMutex
	results  []entry
	scanned  time.Time
	scanning bool
}

func runServer(addr string, refresh time.Duration) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := &server{}
	go func() {
		for {
			s.rescan()
			time.Sleep(refresh)
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handlePage)
	mux.HandleFunc("GET /api/prices", s.handleJSON)
	log.Printf("serving on http://%s (rescan every %s)", ln.Addr(), refresh)
	return http.Serve(ln, mux)
}

func (s *server) rescan() {
	s.mu.Lock()
	s.scanning = true
	s.mu.Unlock()

	start := time.Now()
	results := scan()
	log.Printf("scanned %d products in %s, %d failed", len(results), time.Since(start).Round(time.Second), failures(results))

	s.mu.Lock()
	s.results, s.scanned, s.scanning = results, time.Now(), false
	s.mu.Unlock()
}

type row struct {
	Site, Description, URL, Price, PerGram, Currency, Error string
	Best                                                    bool
}

type group struct {
	Name string
	Rows []row
}

type pageData struct {
	Groups   []group
	Failed   []row
	Scanned  string
	Scanning bool
}

func (s *server) handlePage(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	results, scanned, scanning := s.results, s.scanned, s.scanning
	s.mu.RUnlock()

	data := pageData{Scanning: scanning}
	if !scanned.IsZero() {
		data.Scanned = scanned.Format("2006-01-02 15:04")
	}
	byName := map[string]*group{}
	for _, name := range append(categories, otherCategory) {
		data.Groups = append(data.Groups, group{Name: name})
	}
	for i := range data.Groups {
		byName[data.Groups[i].Name] = &data.Groups[i]
	}
	// Results are already sorted by price per gram.
	for _, res := range results {
		rw := row{Site: res.Site, Description: res.Description, URL: res.URL}
		if res.Err != nil {
			rw.Error = res.Err.Error()
			data.Failed = append(data.Failed, rw)
			continue
		}
		rw.Price, rw.PerGram, rw.Currency = res.Price.StringFixed(2), res.EURPerGram.StringFixed(2), res.Currency
		g := byName[categoryOf(res.Description)]
		rw.Best = len(g.Rows) == 0
		g.Rows = append(g.Rows, rw)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Execute(w, data); err != nil {
		log.Printf("render: %v", err)
	}
}

func (s *server) handleJSON(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	results, scanned := s.results, s.scanned
	s.mu.RUnlock()

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
	items := make([]item, 0, len(results))
	for _, res := range results {
		it := item{Site: res.Site, Category: categoryOf(res.Description), Description: res.Description, URL: res.URL}
		if res.Err != nil {
			it.Error = res.Err.Error()
		} else {
			it.Price, _ = res.Price.Float64()
			it.EURPerGram, _ = res.EURPerGram.Round(4).Float64()
			it.Currency = res.Currency
		}
		items = append(items, it)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"scanned": scanned, "prices": items})
}
