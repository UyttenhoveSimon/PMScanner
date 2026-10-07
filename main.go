package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"pmscanner/scraper"
)

func main() {
	results := scraper.Scrape(products)

	// Cheapest price per gram first; failures last.
	sort.SliceStable(results, func(i, j int) bool {
		a, b := results[i], results[j]
		if (a.Err == nil) != (b.Err == nil) {
			return a.Err == nil
		}
		return a.PricePerGram().LessThan(b.PricePerGram())
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SITE\tPRODUCT\tPRICE\tPER GRAM\tURL")
	failed := 0
	for _, r := range results {
		if r.Err != nil {
			failed++
			fmt.Fprintf(w, "%s\t%s\terror: %v\t\t%s\n", r.Site, r.Description, r.Err, r.URL)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s %s\t%s\t%s\n", r.Site, r.Description,
			r.Price.StringFixed(2), r.Currency, r.PricePerGram().StringFixed(2), r.URL)
	}
	w.Flush()

	if failed == len(results) {
		os.Exit(1)
	}
}
