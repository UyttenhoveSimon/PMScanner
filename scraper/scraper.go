// Package scraper fetches product pages and extracts their price.
//
// Prices are read from structured data (schema.org JSON-LD, microdata or
// Open Graph meta tags) rather than CSS classes, which shops tend to rename.
// Pages behind a bot challenge can be loaded in a headless browser instead.
package scraper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/shopspring/decimal"
)

// Source is a kind of price data found in a page.
type Source int

const (
	CSS       Source = iota // Product.Selectors
	JSONLD                  // <script type="application/ld+json">
	Microdata               // itemprop="price"
	Meta                    // <meta property="product:price:amount">
	numSources
)

// DefaultSources is the order in which price sources are tried.
var DefaultSources = []Source{CSS, JSONLD, Microdata, Meta}

// Product is a product page to scrape.
type Product struct {
	Site        string
	Country     string // ISO 3166-1 alpha-2 code of the shop's country
	Description string
	URL         string
	WeightGrams decimal.Decimal // fine metal content, so prices compare per gram of pure metal
	// Sources overrides DefaultSources, for shops whose preferred source
	// holds the wrong price (e.g. a buy-back price in JSON-LD).
	Sources []Source
	// Selectors are CSS selectors for the price element, tried in order, for
	// shops without structured price data. The element text is the price.
	Selectors []string
	// VAT says how the price is taxed; empty means the default for the metal.
	VAT VAT
	// Currency is used when the page gives a price without a currency.
	Currency string
	// Cookies are sent with the request, e.g. to pick the shop's currency.
	Cookies map[string]string
	// RateGroup makes shops that share a backend, and so a rate limit,
	// share one request limit; by default each host has its own.
	RateGroup string
	// Browser loads the page in a headless browser, for shops that block
	// plain HTTP clients with a JavaScript challenge.
	Browser bool
}

// VAT is how value-added tax applies to a price.
type VAT string

const (
	VATExempt   VAT = "exempt"   // investment gold
	VATIncluded VAT = "included" // VAT is in the price
	VATMargin   VAT = "margin"   // margin scheme: VAT on the dealer's margin, in the price
	VATExcluded VAT = "excluded" // VAT is not in the price but due on purchase
	VATFree     VAT = "free"     // no VAT is due, e.g. metal kept in a vault
)

// Result is the outcome of scraping one Product.
type Result struct {
	Product
	Price    decimal.Decimal
	Currency string
	Err      error
}

// PricePerGram returns the price divided by the product weight.
func (r Result) PricePerGram() decimal.Decimal {
	if r.WeightGrams.IsZero() {
		return decimal.Zero
	}
	return r.Price.Div(r.WeightGrams)
}

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0 Safari/537.36"

var errNoPrice = errors.New("no price found")

const (
	// maxRequests caps plain HTTP requests in flight across all shops.
	maxRequests = 32
	// maxRequestsPerHost keeps the load on each shop low, so a slow shop
	// does not hold up the others and shops are less likely to block us.
	maxRequestsPerHost = 2
)

// hostOf returns the URL's host without a "www." prefix, or the URL itself
// if it cannot be parsed.
func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Hostname() == "" {
		return rawURL
	}
	return strings.TrimPrefix(u.Hostname(), "www.")
}

// hostState limits concurrent requests to one shop and remembers when it
// timed out, so its other products fail fast instead of waiting too.
type hostState struct {
	slots chan struct{}
	mu    sync.Mutex
	err   error
}

func (h *hostState) markUnreachable(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.err == nil {
		h.err = fmt.Errorf("skipped, shop timed out: %w", err)
	}
}

func (h *hostState) unreachable() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.err
}

// Scrape fetches every product concurrently and returns results in input order.
func Scrape(products []Product) []Result {
	results := make([]Result, len(products))
	var browserIdx []int
	var wg sync.WaitGroup
	global := make(chan struct{}, maxRequests)
	hosts := map[string]*hostState{}
	for i, p := range products {
		results[i].Product = p
		if p.Browser {
			browserIdx = append(browserIdx, i)
			continue
		}
		key := p.RateGroup
		if key == "" {
			key = hostOf(p.URL)
		}
		host := hosts[key]
		if host == nil {
			host = &hostState{slots: make(chan struct{}, maxRequestsPerHost)}
			hosts[key] = host
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Take the host slot first so a busy host never holds a global one.
			host.slots <- struct{}{}
			defer func() { <-host.slots }()
			if err := host.unreachable(); err != nil {
				results[i].Err = err
				return
			}
			global <- struct{}{}
			defer func() { <-global }()
			r := &results[i]
			r.Price, r.Currency, r.Err = scrapeHTTP(p)
			var netErr net.Error
			if errors.As(r.Err, &netErr) && netErr.Timeout() {
				host.markUnreachable(r.Err)
			}
		}()
	}
	if len(browserIdx) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scrapeBrowser(results, browserIdx)
		}()
	}
	wg.Wait()
	return results
}

var client = &http.Client{Timeout: 30 * time.Second}

// maxRetries is how many times a request answered with HTTP 429 (too many
// requests) is retried, after the delay the shop asks for.
const maxRetries = 2

func scrapeHTTP(p Product) (decimal.Decimal, string, error) {
	for attempt := 0; ; attempt++ {
		body, retryAfter, err := fetch(p)
		if retryAfter > 0 && attempt < maxRetries {
			time.Sleep(retryAfter)
			continue
		}
		if err != nil {
			return decimal.Zero, "", err
		}
		return extract(body, p)
	}
}

// fetch gets a product page. When the shop rate-limits the request, it also
// returns how long to wait before trying again.
func fetch(p Product) (body []byte, retryAfter time.Duration, err error) {
	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "en,fr;q=0.8,de;q=0.6")
	for name, value := range p.Cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, retryDelay(resp.Header.Get("Retry-After")), fmt.Errorf("fetch: HTTP %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("fetch: HTTP %d", resp.StatusCode)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch: %w", err)
	}
	return body, 0, nil
}

// retryDelay reads a Retry-After header given in seconds, defaulting to a
// few seconds and capped so one shop cannot stall the scan.
func retryDelay(header string) time.Duration {
	const fallback, limit = 3 * time.Second, 15 * time.Second
	secs, err := strconv.Atoi(strings.TrimSpace(header))
	if err != nil || secs <= 0 {
		return fallback
	}
	return min(time.Duration(secs)*time.Second, limit)
}

// extract finds the product price in a page, trying sources in order.
func extract(html []byte, p Product) (decimal.Decimal, string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return decimal.Zero, "", fmt.Errorf("parse: %w", err)
	}
	sources := p.Sources
	if len(sources) == 0 {
		sources = DefaultSources
	}
	for _, src := range sources {
		if price, currency, ok := extractFrom(doc, src, p); ok {
			currency = strings.TrimSpace(currency)
			if currency == "" {
				currency = p.Currency
			}
			return price, currency, nil
		}
	}
	return decimal.Zero, "", errNoPrice
}

func extractFrom(doc *goquery.Document, src Source, p Product) (price decimal.Decimal, currency string, ok bool) {
	switch src {
	case CSS:
		for _, sel := range p.Selectors {
			doc.Find(sel).EachWithBreak(func(_ int, s *goquery.Selection) bool {
				text := s.Text()
				price, ok = parsePositive(text)
				currency = currencyIn(text)
				return !ok
			})
			if ok {
				return
			}
		}
	case JSONLD:
		doc.Find(`script[type="application/ld+json"]`).EachWithBreak(func(_ int, s *goquery.Selection) bool {
			price, currency, ok = fromJSONLD([]byte(s.Text()))
			return !ok
		})
	case Microdata:
		doc.Find(`[itemprop="price"]`).EachWithBreak(func(_ int, s *goquery.Selection) bool {
			raw := s.AttrOr("content", s.Text())
			if price, ok = parsePositive(raw); ok {
				cur := s.Closest(`[itemscope]`).Find(`[itemprop="priceCurrency"]`)
				currency = cur.AttrOr("content", cur.Text())
				if currency == "" {
					currency = currencyIn(s.Text())
				}
			}
			return !ok
		})
	case Meta:
		s := doc.Find(`meta[property="product:price:amount"], meta[property="og:price:amount"]`).First()
		if price, ok = parsePositive(s.AttrOr("content", "")); ok {
			currency = doc.Find(`meta[property$="price:currency"]`).AttrOr("content", "")
		}
	}
	return
}

// currencyIn guesses the ISO currency code from a formatted price.
func currencyIn(s string) string {
	switch {
	case strings.Contains(s, "€"), strings.Contains(s, "EUR"):
		return "EUR"
	case strings.Contains(s, "CHF"):
		return "CHF"
	case strings.Contains(s, "£"), strings.Contains(s, "GBP"):
		return "GBP"
	case strings.Contains(s, "$"), strings.Contains(s, "USD"):
		return "USD"
	case strings.Contains(s, "₽"), strings.Contains(s, "RUB"), strings.Contains(strings.ToLower(s), "руб"):
		return "RUB"
	case strings.Contains(s, "₸"), strings.Contains(s, "KZT"):
		return "KZT"
	case strings.Contains(s, "₼"), strings.Contains(s, "AZN"):
		return "AZN"
	case strings.Contains(s, "₺"), strings.Contains(s, "TRY"):
		return "TRY"
	}
	return ""
}

func parsePositive(s string) (decimal.Decimal, bool) {
	p, err := parsePrice(s)
	return p, err == nil && p.IsPositive()
}

// fromJSONLD finds the first schema.org Offer price in a JSON-LD document.
func fromJSONLD(data []byte) (decimal.Decimal, string, bool) {
	// Some shops emit raw control characters (tabs, newlines) inside JSON
	// strings, which is invalid JSON; whitespace is harmless outside strings.
	data = bytes.Map(func(r rune) rune {
		if r < 0x20 {
			return ' '
		}
		return r
	}, data)
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return decimal.Zero, "", false
	}
	return findOffer(doc)
}

func findOffer(node any) (decimal.Decimal, string, bool) {
	switch v := node.(type) {
	case []any:
		for _, n := range v {
			if p, c, ok := findOffer(n); ok {
				return p, c, ok
			}
		}
	case map[string]any:
		currency, _ := v["priceCurrency"].(string)
		if p, ok := priceField(v, "price"); ok {
			return p, currency, true
		}
		// Nested offers come before an AggregateOffer's lowPrice, which is
		// often a bulk-tier price rather than the single-unit one.
		for _, key := range []string{"offers", "priceSpecification", "@graph", "mainEntity"} {
			if p, c, ok := findOffer(v[key]); ok {
				if c == "" {
					c = currency
				}
				return p, c, ok
			}
		}
		if p, ok := priceField(v, "lowPrice"); ok {
			return p, currency, true
		}
	}
	return decimal.Zero, "", false
}

func priceField(v map[string]any, key string) (decimal.Decimal, bool) {
	raw, ok := v[key]
	if !ok || raw == nil {
		return decimal.Zero, false
	}
	if f, isFloat := raw.(float64); isFloat {
		p := decimal.NewFromFloat(f)
		return p, p.IsPositive()
	}
	return parsePositive(fmt.Sprint(raw))
}

var nonNumeric = regexp.MustCompile(`[^\d.,]`)

// parsePrice parses prices such as "119921.53", "119 921,53 €" or "1,234.50".
func parsePrice(s string) (decimal.Decimal, error) {
	s = nonNumeric.ReplaceAllString(s, "")
	lastDot, lastComma := strings.LastIndex(s, "."), strings.LastIndex(s, ",")
	switch {
	case lastComma > lastDot:
		// Comma is the decimal separator, unless it groups thousands ("1,234").
		if lastDot == -1 && strings.Count(s, ",") == 1 && len(s)-lastComma-1 == 3 {
			s = strings.ReplaceAll(s, ",", "")
		} else {
			s = strings.ReplaceAll(s, ".", "")
			s = strings.Replace(s, ",", ".", 1)
		}
	case lastDot > lastComma:
		s = strings.ReplaceAll(s, ",", "")
		// Several dots means they group thousands ("1.234.567").
		if strings.Count(s, ".") > 1 {
			s = strings.ReplaceAll(s, ".", "")
		}
	}
	if s == "" {
		return decimal.Zero, errors.New("empty price")
	}
	return decimal.NewFromString(s)
}
