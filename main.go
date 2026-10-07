package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"pmscanner/scraper"
)

func main() {
	serve := flag.String("serve", "", "serve a web page on this address (e.g. :8080) instead of printing a table")
	refresh := flag.Duration("refresh", 30*time.Minute, "how often the web page rescans prices")
	flag.Parse()

	if *serve != "" {
		log.Fatal(runServer(*serve, *refresh))
	}

	results := scan()
	printTable(results)
	if failures(results) == len(results) {
		os.Exit(1)
	}
}

// scan scrapes all products, cheapest price per gram first and failures last.
func scan() []scraper.Result {
	results := scraper.Scrape(products)
	sort.SliceStable(results, func(i, j int) bool {
		a, b := results[i], results[j]
		if (a.Err == nil) != (b.Err == nil) {
			return a.Err == nil
		}
		return a.PricePerGram().LessThan(b.PricePerGram())
	})
	return results
}

func failures(results []scraper.Result) int {
	n := 0
	for _, r := range results {
		if r.Err != nil {
			n++
		}
	}
	return n
}

func printTable(results []scraper.Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SITE\tPRODUCT\tPRICE\tPER GRAM\tURL")
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(w, "%s\t%s\terror: %v\t\t%s\n", r.Site, r.Description, r.Err, r.URL)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s %s\t%s\t%s\n", r.Site, r.Description,
			r.Price.StringFixed(2), r.Currency, r.PricePerGram().StringFixed(2), r.URL)
	}
	w.Flush()
}
