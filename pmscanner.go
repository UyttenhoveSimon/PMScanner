package main

import (
	"github.com/gocolly/colly/v2"
)

func main() {

	// Instantiate default collector
	c := colly.NewCollector()

	// For each go in websites/, we use its tailored method to gather prices

	c.Visit("")
}
