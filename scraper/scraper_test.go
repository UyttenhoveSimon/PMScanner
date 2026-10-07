package scraper

import "testing"

func TestParsePrice(t *testing.T) {
	cases := map[string]string{
		"119921.53":      "119921.53",
		"119 921,53 €":   "119921.53",
		"119.921,53":     "119921.53",
		"1,234.50":       "1234.5",
		"1,234":          "1234",
		"2.345,6":        "2345.6",
		"1.234.567":      "1234567",
		"€ 3 512,00 TTC": "3512",
		"42":             "42",
	}
	for in, want := range cases {
		got, err := parsePrice(in)
		if err != nil || got.String() != want {
			t.Errorf("parsePrice(%q) = %v, %v; want %s", in, got, err, want)
		}
	}
}

func TestFromJSONLD(t *testing.T) {
	doc := `{"@graph":[{"@type":"Product","offers":[{"@type":"Offer","price":"2 345,60","priceCurrency":"EUR"}]}]}`
	p, c, ok := fromJSONLD([]byte(doc))
	if !ok || p.String() != "2345.6" || c != "EUR" {
		t.Errorf("got %v %q %v", p, c, ok)
	}
}
