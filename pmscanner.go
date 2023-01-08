package main

import (
	"fmt"
	. "pmscanner/websites"

	"github.com/gocolly/colly/v2"
)

func main() {

	// Instantiate default collector
	c := colly.NewCollector()

	orWebsite := Or{}
	collectedItems := orWebsite.Collect(c)

	fmt.Printf("%v\n", collectedItems)

	// For each go in websites/, we use its tailored method to gather prices
}
