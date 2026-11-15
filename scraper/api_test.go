package scraper

import (
	"net/url"
	"strings"
	"testing"
)

func TestWooCommerceURL(t *testing.T) {
	p := Product{URL: "https://centerzlata.com/izdelek/krugerrand-1-unca-zlati-kovanec/", API: WooCommerce}
	got, err := apiURL(p)
	want := "https://centerzlata.com/wp-json/wc/store/v1/products?slug=krugerrand-1-unca-zlati-kovanec"
	if err != nil || got != want {
		t.Errorf("apiURL = %q, %v; want %q", got, err, want)
	}
}

func TestFromWooCommerce(t *testing.T) {
	body := []byte(`[{"id":429,"slug":"krugerrand","prices":{"price":"360964","currency_code":"EUR","currency_minor_unit":2}}]`)
	price, currency, err := fromWooCommerce(body, "krugerrand")
	if err != nil || price.String() != "3609.64" || currency != "EUR" {
		t.Errorf("got %s %s %v, want 3609.64 EUR", price, currency, err)
	}
	if _, _, err := fromWooCommerce([]byte(`[]`), "krugerrand"); err != errNoPrice {
		t.Errorf("empty list: got %v, want errNoPrice", err)
	}
	// Some shops ignore the slug filter and answer with another product.
	if _, _, err := fromWooCommerce(body, "maple-leaf"); err != errNoPrice {
		t.Errorf("other product: got %v, want errNoPrice", err)
	}
}

func TestMagentoURL(t *testing.T) {
	p := Product{URL: "https://www.example.com/nl/silver-bar-1kg.html", API: Magento}
	got, err := apiURL(p)
	if err != nil || !strings.HasPrefix(got, "https://www.example.com/graphql?query=") ||
		!strings.Contains(got, url.QueryEscape(`{eq:"silver-bar-1kg"}`)) {
		t.Errorf("apiURL = %q, %v", got, err)
	}
}

func TestFromMagento(t *testing.T) {
	body := []byte(`{"data":{"products":{"items":[{"url_key":"bar","price_range":{"minimum_price":{"final_price":{"value":1844.67,"currency":"EUR"}}}}]}}}`)
	price, currency, err := fromMagento(body, "bar")
	if err != nil || price.String() != "1844.67" || currency != "EUR" {
		t.Errorf("got %s %s %v, want 1844.67 EUR", price, currency, err)
	}
	if _, _, err := fromMagento(body, "other"); err != errNoPrice {
		t.Errorf("other product: got %v, want errNoPrice", err)
	}
	broken := []byte(`{"errors":[{"message":"Internal server error"}],"data":{"products":null}}`)
	if _, _, err := fromMagento(broken, "bar"); err == nil {
		t.Error("graphql error: got nil error")
	}
}
