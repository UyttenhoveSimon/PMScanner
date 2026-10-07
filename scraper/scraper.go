// Package scraper fetches product pages and extracts their price.
//
// Prices are read from structured data (schema.org JSON-LD, microdata or
// Open Graph meta tags) rather than CSS classes, which shops tend to rename.
package scraper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
)

// Source is a kind of structured price data found in a page.
type Source int

const (
	JSONLD    Source = iota // <script type="application/ld+json">
	Microdata               // itemprop="price"
	Meta                    // <meta property="product:price:amount">
	numSources
)

// DefaultSources is the order in which price sources are tried.
var DefaultSources = []Source{JSONLD, Microdata, Meta}

// Product is a product page to scrape.
type Product struct {
	Site        string
	Description string
	URL         string
	WeightGrams decimal.Decimal
	// Sources overrides DefaultSources, for shops whose preferred source
	// holds the wrong price (e.g. a buy-back price in JSON-LD).
	Sources []Source
}

type candidate struct {
	price    decimal.Decimal
	currency string
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

// Scrape fetches every product concurrently and returns results in input order.
func Scrape(products []Product) []Result {
	results := make([]Result, len(products))
	for i, p := range products {
		results[i] = Result{Product: p, Err: errors.New("no price found")}
	}

	c := colly.NewCollector(colly.Async(true), colly.UserAgent(userAgent))
	c.SetRequestTimeout(30 * time.Second)
	_ = c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})

	index := func(r *colly.Response) int { return r.Ctx.GetAny("index").(int) }
	candidates := func(r *colly.Response) *[numSources]*candidate {
		return r.Ctx.GetAny("candidates").(*[numSources]*candidate)
	}
	// Keeps the first price seen for each source.
	add := func(r *colly.Response, src Source, price decimal.Decimal, currency string) {
		if cs := candidates(r); cs[src] == nil && price.IsPositive() {
			cs[src] = &candidate{price, strings.TrimSpace(currency)}
		}
	}

	c.OnHTML(`script[type="application/ld+json"]`, func(e *colly.HTMLElement) {
		if price, currency, ok := fromJSONLD([]byte(e.Text)); ok {
			add(e.Response, JSONLD, price, currency)
		}
	})
	c.OnHTML(`[itemprop="price"]`, func(e *colly.HTMLElement) {
		raw := e.Attr("content")
		if raw == "" {
			raw = e.Text
		}
		if price, err := parsePrice(raw); err == nil {
			cur := e.DOM.Closest(`[itemscope]`).Find(`[itemprop="priceCurrency"]`)
			currency := cur.AttrOr("content", cur.Text())
			add(e.Response, Microdata, price, currency)
		}
	})
	c.OnHTML(`meta[property="product:price:amount"], meta[property="og:price:amount"]`, func(e *colly.HTMLElement) {
		if price, err := parsePrice(e.Attr("content")); err == nil {
			currency := e.DOM.Parent().Find(`meta[property$="price:currency"]`).AttrOr("content", "")
			add(e.Response, Meta, price, currency)
		}
	})
	c.OnScraped(func(r *colly.Response) {
		res := &results[index(r)]
		sources := res.Sources
		if len(sources) == 0 {
			sources = DefaultSources
		}
		for _, src := range sources {
			if cand := candidates(r)[src]; cand != nil {
				res.Price, res.Currency, res.Err = cand.price, cand.currency, nil
				return
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		results[index(r)].Err = fmt.Errorf("fetch: %w (HTTP %d)", err, r.StatusCode)
	})

	for i, p := range products {
		ctx := colly.NewContext()
		ctx.Put("index", i)
		ctx.Put("candidates", &[numSources]*candidate{})
		if err := c.Request("GET", p.URL, nil, ctx, nil); err != nil {
			results[i].Err = fmt.Errorf("fetch: %w", err)
		}
	}
	c.Wait()
	return results
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
		return decimal.NewFromFloat(f), true
	}
	p, err := parsePrice(fmt.Sprint(raw))
	return p, err == nil && p.IsPositive()
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
