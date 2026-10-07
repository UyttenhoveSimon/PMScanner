package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"pmscanner/store"
)

func TestHistoryRoundTrip(t *testing.T) {
	dir := t.TempDir()
	src, err := store.Open(filepath.Join(dir, "src.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()
	at := time.Date(2026, 10, 7, 5, 0, 0, 0, time.UTC)
	prices := []store.Price{{Site: "a", Country: "BE", VAT: "exempt", Category: "Gold bar 1kg", Description: "Gold bar 1kg",
		URL: "https://a", Price: decimal.RequireFromString("119000.5"), Currency: "EUR", EURPerGram: decimal.RequireFromString("119.0005")}}
	if _, err := src.Save(at, prices); err != nil {
		t.Fatal(err)
	}
	hist := filepath.Join(dir, "scans")
	if err := exportHistory(hist, src); err != nil {
		t.Fatal(err)
	}
	if err := exportHistory(hist, src); err != nil { // idempotent
		t.Fatal(err)
	}

	dst, err := store.Open(filepath.Join(dir, "dst.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer dst.Close()
	for range 2 { // importing twice must not duplicate
		if err := importHistory(hist, dst); err != nil {
			t.Fatal(err)
		}
	}
	scans, err := dst.Scans(10)
	if err != nil || len(scans) != 1 {
		t.Fatalf("scans = %v, %v; want 1", scans, err)
	}
	when, got, err := dst.Prices(scans[0].ID)
	if err != nil || !when.Equal(at) || len(got) != 1 || !got[0].Price.Equal(prices[0].Price) || got[0].VAT != "exempt" {
		t.Errorf("got %v %+v %v", when, got, err)
	}
}
