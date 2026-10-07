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
	"net/http"
	"regexp"
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
	Description string
	URL         string
	WeightGrams decimal.Decimal
	// Sources overrides DefaultSources, for shops whose preferred source
	// holds the wrong price (e.g. a buy-back price in JSON-LD).
	Sources []Source
	// Selectors are CSS selectors for the price element, tried in order, for
	// shops without structured price data. The element text is the price.
	Selectors []string
	// Cookies are sent with the request, e.g. to pick the shop's currency.
	Cookies map[string]string
	// Browser loads the page in a headless browser, for shops that block
	// plain HTTP clients with a JavaScript challenge.
	Browser bool
}

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

// Scrape fetches every product concurrently and returns results in input order.
func Scrape(products []Product) []Result {
	results := make([]Result, len(products))
	var browserIdx []int
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	for i, p := range products {
		results[i].Product = p
		if p.Browser {
			browserIdx = append(browserIdx, i)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i].Price, results[i].Currency, results[i].Err = scrapeHTTP(p)
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

func scrapeHTTP(p Product) (decimal.Decimal, string, error) {
	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return decimal.Zero, "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "en,fr;q=0.8,de;q=0.6")
	for name, value := range p.Cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Zero, "", fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, "", fmt.Errorf("fetch: HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, "", fmt.Errorf("fetch: %w", err)
	}
	return extract(body, p)
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
			return price, strings.TrimSpace(currency), nil
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
