package scraper

import "testing"

func TestWooCommerceURL(t *testing.T) {
	p := Product{URL: "https://centerzlata.com/izdelek/krugerrand-1-unca-zlati-kovanec/", API: WooCommerce}
	got, err := apiURL(p)
	want := "https://centerzlata.com/wp-json/wc/store/v1/products?slug=krugerrand-1-unca-zlati-kovanec"
	if err != nil || got != want {
		t.Errorf("apiURL = %q, %v; want %q", got, err, want)
	}
}

func TestFromWooCommerce(t *testing.T) {
	body := []byte(`[{"id":429,"prices":{"price":"360964","currency_code":"EUR","currency_minor_unit":2}}]`)
	price, currency, err := fromWooCommerce(body)
	if err != nil || price.String() != "3609.64" || currency != "EUR" {
		t.Errorf("got %s %s %v, want 3609.64 EUR", price, currency, err)
	}
	if _, _, err := fromWooCommerce([]byte(`[]`)); err != errNoPrice {
		t.Errorf("empty list: got %v, want errNoPrice", err)
	}
}
