package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"pmscanner/fx"
	"pmscanner/store"
)

// staticLinks points to the files written by exportSite. Pages of past scans
// live one directory down, so links are relative to the page's depth.
type staticLinks struct{ root string }

func (l staticLinks) Scan(id int64, latest bool) string {
	if latest {
		return l.root + "index.html"
	}
	return l.root + "scans/" + strconv.FormatInt(id, 10) + ".html"
}
func (l staticLinks) JSON(id int64) string {
	return l.root + "api/prices/" + strconv.FormatInt(id, 10) + ".json"
}
func (l staticLinks) Latest() string { return l.root + "index.html" }

// exportSite writes a static copy of the web page to dir, for hosting
// without a server (e.g. GitHub Pages):
//
//	index.html             latest scan
//	scans/ID.html          every saved scan
//	api/scans.json         list of scans
//	api/prices.json        latest scan's prices
//	api/prices/ID.json     every scan's prices
//	sitemap.xml, llms.txt  for search engines and language models
func exportSite(dir string, db *store.Store) error {
	scans, err := db.Scans(scanListLimit)
	if err != nil {
		return err
	}
	if len(scans) == 0 {
		return fmt.Errorf("no scan to export; run a scan first")
	}
	rates, origins, err := fx.FetchRates()
	if err != nil {
		log.Printf("exchange rates: %v (the page will only offer EUR)", err)
	}

	for _, sub := range []string{"scans", "api/prices"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			return err
		}
	}
	for i, sc := range scans {
		name := filepath.Join("scans", strconv.FormatInt(sc.ID, 10)+".html")
		if err := writePage(dir, name, db, sc.ID, staticLinks{root: "../"}, rates, origins); err != nil {
			return err
		}
		prices, err := pricesJSON(db, sc.ID)
		if err != nil {
			return err
		}
		if err := writeJSONFile(filepath.Join(dir, "api", "prices", strconv.FormatInt(sc.ID, 10)+".json"), prices); err != nil {
			return err
		}
		if i == 0 {
			if err := writePage(dir, "index.html", db, sc.ID, staticLinks{}, rates, origins); err != nil {
				return err
			}
			if err := writeJSONFile(filepath.Join(dir, "api", "prices.json"), prices); err != nil {
				return err
			}
		}
	}
	if err := writeJSONFile(filepath.Join(dir, "api", "scans.json"), scansJSON(scans)); err != nil {
		return err
	}
	if err := writeCrawlerFiles(dir, scans[0].At); err != nil {
		return err
	}
	log.Printf("exported %d scans to %s", len(scans), dir)
	return nil
}

func writePage(dir, name string, db *store.Store, id int64, l links, rates fx.Rates, origins fx.Origins) error {
	data, err := buildPage(db, id, l, rates, origins)
	if err != nil {
		return fmt.Errorf("scan %d: %w", id, err)
	}
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	if data.IsLatest {
		data.Canonical = siteURL
	}
	if err := renderPage(f, data); err != nil {
		f.Close()
		return fmt.Errorf("render %s: %w", name, err)
	}
	return f.Close()
}

func writeJSONFile(path string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// siteURL is where the Pages workflow publishes the export. Only the latest
// scan is indexed: past scans are the same page with older prices.
const siteURL = "https://uyttenhovesimon.github.io/PMScanner/"

func writeCrawlerFiles(dir string, scanned time.Time) error {
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>` + siteURL + `</loc><lastmod>` + scanned.UTC().Format("2006-01-02") + `</lastmod><changefreq>daily</changefreq></url>
</urlset>
`
	llms := `# PM Scanner

> Gold, silver, platinum and palladium bar and coin prices from European
> bullion dealers, ranked by price per gram with VAT included. Scanned daily.

The page lists one table per product type (e.g. "Gold bar 1kg", "Gold
Krugerrand 1oz", "Silver bar 1kg"), cheapest first. The same data is
available as JSON.

## Data

- [Latest prices](` + siteURL + `api/prices.json): one entry per product: site, country (ISO code), category, description, url, price, currency, EUR per gram, VAT treatment
- [Scans](` + siteURL + `api/scans.json): every saved scan with its id and date; prices per scan at api/prices/ID.json

## Source

- [GitHub repository](https://github.com/UyttenhoveSimon/PMScanner): the scraper, the list of shops, VAT rules by country
`
	if err := os.WriteFile(filepath.Join(dir, "sitemap.xml"), []byte(sitemap), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "llms.txt"), []byte(llms), 0o644)
}
