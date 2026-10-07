package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/shopspring/decimal"

	"pmscanner/fx"
	"pmscanner/scraper"
	"pmscanner/store"
)

func main() {
	serve := flag.String("serve", "", "serve a web page on this address (e.g. :8080) instead of printing a table")
	refresh := flag.Duration("refresh", 24*time.Hour, "how often the web server rescans prices (0 disables automatic scans)")
	dbPath := flag.String("db", "pmscanner.db", "SQLite database where every scan is saved")
	flag.Parse()

	db, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if *serve != "" {
		if err := runServer(*serve, *refresh, db); err != nil {
			db.Close()
			log.Fatal(err)
		}
		return
	}

	start := time.Now()
	entries := scan()
	printTable(entries)
	if _, err := db.Save(start, toPrices(entries)); err != nil {
		log.Printf("save scan: %v", err)
	}
	if failures(entries) == len(entries) {
		db.Close()
		os.Exit(1)
	}
}

// entry is a scraped product with its price per gram converted to euros.
type entry struct {
	scraper.Result
	EURPerGram decimal.Decimal
}

// scan scrapes all products, cheapest price per gram first and failures last.
func scan() []entry {
	rates, ratesErr := fx.FetchECB()
	if ratesErr != nil {
		log.Printf("prices in other currencies than EUR will fail: %v", ratesErr)
	}
	results := scraper.Scrape(products)
	entries := make([]entry, len(results))
	for i, r := range results {
		e := entry{Result: r}
		if r.Err == nil {
			switch eur, err := rates.ToEUR(r.PricePerGram(), r.Currency); {
			case err != nil && ratesErr != nil:
				e.Err = ratesErr
			case err != nil:
				e.Err = err
			default:
				e.EURPerGram = eur
			}
		}
		entries[i] = e
	}
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if (a.Err == nil) != (b.Err == nil) {
			return a.Err == nil
		}
		return a.EURPerGram.LessThan(b.EURPerGram)
	})
	return entries
}

// toPrices converts scan entries to database rows, keeping their order.
func toPrices(entries []entry) []store.Price {
	prices := make([]store.Price, len(entries))
	for i, e := range entries {
		p := store.Price{
			Site: e.Site, Category: categoryOf(e.Description), Description: e.Description, URL: e.URL,
			Price: e.Price, Currency: e.Currency, EURPerGram: e.EURPerGram,
		}
		if e.Err != nil {
			p.Error = e.Err.Error()
		}
		prices[i] = p
	}
	return prices
}

func failures(entries []entry) int {
	n := 0
	for _, e := range entries {
		if e.Err != nil {
			n++
		}
	}
	return n
}

func printTable(entries []entry) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SITE\tPRODUCT\tPRICE\tEUR/GRAM\tURL")
	for _, r := range entries {
		if r.Err != nil {
			fmt.Fprintf(w, "%s\t%s\terror: %v\t\t%s\n", r.Site, r.Description, r.Err, r.URL)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s %s\t%s\t%s\n", r.Site, r.Description,
			r.Price.StringFixed(2), r.Currency, r.EURPerGram.StringFixed(2), r.URL)
	}
	w.Flush()
}
