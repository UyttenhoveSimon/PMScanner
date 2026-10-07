package main

import (
	"github.com/shopspring/decimal"

	"pmscanner/scraper"
)

var (
	kg = decimal.NewFromInt(1000)
)

// products lists the pages to scan. Any shop page exposing schema.org price
// data (JSON-LD, microdata or Open Graph) works without site-specific code.
var products = []scraper.Product{
	{Site: "or.fr", Description: "Gold bar 1kg Valcambi", WeightGrams: kg,
		URL: "https://or.fr/achat/or/lingots/lingot-d-or-1-kilogramme-valcambi-60"},
}
