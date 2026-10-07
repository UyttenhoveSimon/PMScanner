package main

import (
	"testing"

	"github.com/shopspring/decimal"

	"pmscanner/scraper"
)

func TestWithVAT(t *testing.T) {
	two := decimal.NewFromInt(2)
	cases := []struct {
		vat     scraper.VAT
		country string
		want    string
		ok      bool
	}{
		{scraper.VATExcluded, "ES", "2.42", true}, // +21%
		{scraper.VATExcluded, "FI", "2.51", true}, // +25.5%
		{scraper.VATIncluded, "ES", "2", true},
		{scraper.VATFree, "FR", "2", true},
		{scraper.VATExcluded, "XX", "2", false},
	}
	for _, c := range cases {
		got, ok := withVAT(two, c.vat, c.country)
		if got.String() != c.want || ok != c.ok {
			t.Errorf("withVAT(2, %s, %s) = %s, %v; want %s, %v", c.vat, c.country, got, ok, c.want, c.ok)
		}
	}
}

func TestEveryProductCountryHasTaxes(t *testing.T) {
	for _, p := range products {
		if _, ok := taxes[p.Country]; !ok {
			t.Errorf("%s (%s): no tax entry for country %q", p.Site, p.Description, p.Country)
		}
	}
}
