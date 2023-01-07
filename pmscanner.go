package main

import (
	. "pmscanner/websites"

	"github.com/gocolly/colly/v2"
)

func main() {

	// Instantiate default collector
	c := colly.NewCollector()

	orPrices := Or{}

	// For each go in websites/, we use its tailored method to gather prices
}
