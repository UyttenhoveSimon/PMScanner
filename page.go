package main

import (
	_ "embed"
	"html/template"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	"pmscanner/fx"
	"pmscanner/scraper"
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
	"Gold Philharmonic 1oz",
	"Gold Britannia 1oz",
	"Gold Eagle 1oz",
	"Gold Kangaroo 1oz",
	"Gold Sovereign",
	"Gold 20 Mark",
	"Gold Vreneli 20 CHF",
	"Gold Napoleon 20 FF",
	"Silver Maple Leaf 1oz",
	"Silver bar 1kg",
	"Platinum bar 1oz",
	"Platinum coin 1oz",
	"Palladium bar 1oz",
	"Palladium coin 1oz",
	"Copper bar 1kg",
	"Rhodium",
	"Iridium",
	"Ruthenium",
	"Osmium",
}

const otherCategory = "Other"

// metalOf returns a category's metal, its first word ("Gold bar 1kg" → "Gold").
func metalOf(category string) string {
	metal, _, _ := strings.Cut(category, " ")
	return metal
}

func categoryOf(description string) string {
	for _, c := range categories {
		if strings.HasPrefix(description, c) {
			return c
		}
	}
	return otherCategory
}

// productByURL fills in details of prices saved before they were recorded.
var productByURL = func() map[string]scraper.Product {
	m := make(map[string]scraper.Product, len(products))
	for _, p := range products {
		m[p.URL] = p
	}
	return m
}()

// vatOf returns how a product's price is taxed: investment gold is VAT-exempt
// in the EU, UK and Switzerland, and other metals include VAT unless the
// product says otherwise.
func vatOf(p scraper.Product) scraper.VAT {
	if p.VAT != "" {
		return p.VAT
	}
	if metalOf(categoryOf(p.Description)) == "Gold" {
		return scraper.VATExempt
	}
	return scraper.VATIncluded
}

// vatLabels are shown next to prices.
var vatLabels = map[scraper.VAT]string{
	scraper.VATExempt:   "VAT exempt",
	scraper.VATIncluded: "incl. VAT",
	scraper.VATMargin:   "margin VAT",
	scraper.VATExcluded: "excl. VAT",
}

// savedVAT returns a saved price's VAT, working it out for scans saved
// before it was recorded.
func savedVAT(p store.Price) scraper.VAT {
	if p.VAT != "" {
		return scraper.VAT(p.VAT)
	}
	if prod, ok := productByURL[p.URL]; ok {
		return vatOf(prod)
	}
	return vatOf(scraper.Product{Description: p.Description})
}

// scanListLimit caps how many past scans are listed and exported.
const scanListLimit = 500

const timeFormat = "2006-01-02 15:04 MST"

// links builds the URLs a page points to, which differ between the live
// server and the static export.
type links interface {
	Scan(id int64, latest bool) string
	JSON(id int64) string
	Latest() string
}

type row struct {
	Site, Country, Description, URL, Price, PerGram, Currency, Error string
	VAT, VATLabel                                                    string
	Best                                                             bool
}

type group struct {
	Name  string
	Metal string
	Rows  []row
}

type scanOption struct {
	URL      string
	Label    string
	Selected bool
}

type pageData struct {
	Groups     []group
	Rates      map[string]float64 // units per euro, for switching currency in the page
	Currencies []string
	Countries  []string
	Metals     []string // in category order
	Failed     []row
	Scanned    string
	Scanning   bool
	IsLatest   bool
	Scans      []scanOption
	JSONURL    string
	LatestURL  string
}

// buildPage gathers everything the page shows for one scan.
func buildPage(db *store.Store, id int64, l links, rates fx.Rates) (pageData, error) {
	var data pageData
	data.Rates, data.Currencies = currencyMenu(rates)

	scanned, prices, err := db.Prices(id)
	if err != nil {
		return data, err
	}
	scans, err := db.Scans(scanListLimit)
	if err != nil {
		return data, err
	}

	data.Scanned = scanned.Local().Format(timeFormat)
	data.IsLatest = len(scans) > 0 && scans[0].ID == id
	data.JSONURL, data.LatestURL = l.JSON(id), l.Latest()
	for i, sc := range scans {
		label := sc.At.Local().Format(timeFormat)
		if sc.Failed > 0 {
			label += " (" + strconv.Itoa(sc.Failed) + " failed)"
		}
		data.Scans = append(data.Scans, scanOption{URL: l.Scan(sc.ID, i == 0), Label: label, Selected: sc.ID == id})
	}

	byName := map[string]*group{}
	for _, name := range append(categories, otherCategory) {
		data.Groups = append(data.Groups, group{Name: name, Metal: metalOf(name)})
	}
	for i := range data.Groups {
		byName[data.Groups[i].Name] = &data.Groups[i]
	}
	countries := map[string]bool{}
	// Prices are saved sorted by price per gram.
	for _, p := range prices {
		country := p.Country
		if country == "" {
			country = productByURL[p.URL].Country // scans saved before countries were recorded
		}
		rw := row{Site: p.Site, Country: country, Description: p.Description, URL: p.URL}
		if p.Error != "" {
			rw.Error = p.Error
			data.Failed = append(data.Failed, rw)
			continue
		}
		rw.Price, rw.PerGram, rw.Currency = p.Price.StringFixed(2), p.EURPerGram.StringFixed(2), p.Currency
		vat := savedVAT(p)
		rw.VAT, rw.VATLabel = string(vat), vatLabels[vat]
		if country != "" {
			countries[country] = true
		}
		g, ok := byName[p.Category]
		if !ok {
			g = byName[otherCategory]
		}
		rw.Best = len(g.Rows) == 0
		g.Rows = append(g.Rows, rw)
	}
	for c := range countries {
		data.Countries = append(data.Countries, c)
	}
	for _, g := range data.Groups {
		if len(g.Rows) > 0 && !slices.Contains(data.Metals, g.Metal) {
			data.Metals = append(data.Metals, g.Metal)
		}
	}
	slices.Sort(data.Countries)
	return data, nil
}

func renderPage(w io.Writer, data pageData) error {
	return page.Execute(w, data)
}

// preferredCurrencies are listed first in the currency menu.
var preferredCurrencies = []string{"EUR", "USD", "GBP", "CHF"}

// currencyMenu returns the exchange rates for the page, as units per euro,
// and the currency menu. Without rates, only EUR is offered.
func currencyMenu(rates fx.Rates) (map[string]float64, []string) {
	out := map[string]float64{"EUR": 1}
	for cur, r := range rates {
		out[cur], _ = r.Float64()
	}
	names := slices.DeleteFunc(slices.Clone(preferredCurrencies), func(c string) bool {
		_, ok := out[c]
		return !ok
	})
	var others []string
	for cur := range out {
		if !slices.Contains(names, cur) {
			others = append(others, cur)
		}
	}
	slices.Sort(others)
	return out, append(names, others...)
}

type priceJSON struct {
	Site        string  `json:"site"`
	Country     string  `json:"country,omitempty"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	Price       float64 `json:"price,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	EURPerGram  float64 `json:"eurPerGram,omitempty"`
	VAT         string  `json:"vat,omitempty"`
	Error       string  `json:"error,omitempty"`
}

// pricesJSON is the API representation of one scan.
func pricesJSON(db *store.Store, id int64) (any, error) {
	scanned, prices, err := db.Prices(id)
	if err != nil {
		return nil, err
	}
	items := make([]priceJSON, 0, len(prices))
	for _, p := range prices {
		it := priceJSON{Site: p.Site, Country: p.Country, Category: p.Category, Description: p.Description, URL: p.URL, Error: p.Error}
		if it.Country == "" {
			it.Country = productByURL[p.URL].Country
		}
		if p.Error == "" {
			it.Price, _ = p.Price.Float64()
			it.EURPerGram, _ = p.EURPerGram.Round(4).Float64()
			it.Currency = p.Currency
			it.VAT = string(savedVAT(p))
		}
		items = append(items, it)
	}
	return map[string]any{"scan": id, "scanned": scanned, "prices": items}, nil
}

type scanJSON struct {
	ID     int64     `json:"id"`
	At     time.Time `json:"at"`
	Count  int       `json:"count"`
	Failed int       `json:"failed"`
}

// scansJSON is the API representation of the scan list.
func scansJSON(scans []store.Scan) any {
	items := make([]scanJSON, 0, len(scans))
	for _, sc := range scans {
		items = append(items, scanJSON{sc.ID, sc.At, sc.Count, sc.Failed})
	}
	return map[string]any{"scans": items}
}
