package main

import (
	"pmscanner/scraper"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCurrencyMenu(t *testing.T) {
	rates := map[string]decimal.Decimal{
		"USD": decimal.RequireFromString("1.12"),
		"CHF": decimal.RequireFromString("0.93"),
		"JPY": decimal.RequireFromString("178"),
	}
	got, names := currencyMenu(rates, []string{"CHF", "EUR", "CHF", "XYZ", ""})
	want := []string{"CHF", "EUR", "USD"}
	if len(names) != len(want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("names = %v, want %v", names, want)
		}
	}
	if _, ok := got["JPY"]; ok {
		t.Error("JPY should not be offered: no shop uses it")
	}
}

func TestVATFor(t *testing.T) {
	silver := scraper.Product{Description: "Silver bar 1kg", ExportCurrency: "USD"}
	gold := scraper.Product{Description: "Gold bar 1kg", ExportCurrency: "USD"}
	cases := []struct {
		p        scraper.Product
		currency string
		want     scraper.VAT
	}{
		{silver, "CHF", scraper.VATIncluded},
		{silver, "USD", scraper.VATExcluded},
		{gold, "USD", scraper.VATExempt},
		{scraper.Product{Description: "Silver bar 1kg"}, "USD", scraper.VATIncluded},
	}
	for _, c := range cases {
		if got := vatFor(c.p, c.currency); got != c.want {
			t.Errorf("vatFor(%q, %s) = %s, want %s", c.p.Description, c.currency, got, c.want)
		}
	}
}
