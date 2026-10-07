package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestSaveAndLoad(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if _, err := s.Latest(); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Latest on empty db: %v, want ErrNotFound", err)
	}

	at := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	in := []Price{
		{Site: "b", Country: "BE", VAT: "exempt", Category: "Gold bar 1kg", Description: "Gold bar 1kg", URL: "https://b",
			Price: decimal.RequireFromString("119000.50"), Currency: "EUR",
			EURPerGram: decimal.RequireFromString("119.0005")},
		{Site: "a", Category: "Gold bar 1kg", Description: "Gold bar 1kg", URL: "https://a",
			Error: "no price found"},
	}
	first, err := s.Save(at, in)
	if err != nil {
		t.Fatal(err)
	}
	second, err := s.Save(at.Add(time.Hour), in[:1])
	if err != nil {
		t.Fatal(err)
	}

	if latest, err := s.Latest(); err != nil || latest != second {
		t.Errorf("Latest = %d, %v; want %d", latest, err, second)
	}

	scans, err := s.Scans(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(scans) != 2 || scans[0].ID != second || scans[1].Count != 2 || scans[1].Failed != 1 {
		t.Errorf("Scans = %+v", scans)
	}

	when, out, err := s.Prices(first)
	if err != nil {
		t.Fatal(err)
	}
	if !when.Equal(at) {
		t.Errorf("time = %v, want %v", when, at)
	}
	if len(out) != 2 || out[0].Site != "b" || out[0].Country != "BE" || out[0].VAT != "exempt" || !out[0].Price.Equal(in[0].Price) ||
		!out[0].EURPerGram.Equal(in[0].EURPerGram) || out[1].Error != "no price found" {
		t.Errorf("Prices = %+v", out)
	}

	if _, _, err := s.Prices(999); !errors.Is(err, ErrNotFound) {
		t.Errorf("Prices(999) err = %v, want ErrNotFound", err)
	}
}

func TestMigrateAddsCountry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "old.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	// Recreate the prices table as it was before the country column.
	if _, err := s.db.Exec(`DROP TABLE prices; CREATE TABLE prices (
		scan_id INTEGER NOT NULL, position INTEGER NOT NULL, site TEXT NOT NULL,
		category TEXT NOT NULL, description TEXT NOT NULL, url TEXT NOT NULL,
		price TEXT NOT NULL DEFAULT '', currency TEXT NOT NULL DEFAULT '',
		eur_per_gram TEXT NOT NULL DEFAULT '', error TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (scan_id, position))`); err != nil {
		t.Fatal(err)
	}
	s.Close()

	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err := s.Save(time.Now(), []Price{{Site: "a", Country: "FR", Error: "x"}}); err != nil {
		t.Fatalf("save after migration: %v", err)
	}
}
