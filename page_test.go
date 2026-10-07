package main

import (
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
