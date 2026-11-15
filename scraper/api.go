package scraper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/shopspring/decimal"
)

// API is a shop JSON API that gives a product's price without loading the
// page. Shops often leave these open when the page itself is bot-walled.
type API string

const (
	// WooCommerce is the WooCommerce Store API, looked up by the product's
	// slug: the last segment of the page URL. Its prices follow the shop's
	// tax setting, which can leave VAT out even when the page includes it.
	WooCommerce API = "woocommerce"
	// Magento is the Magento (Adobe Commerce) GraphQL API, looked up by the
	// product's url_key: the last segment of the page URL, without ".html".
	Magento API = "magento"
)

// magentoQuery asks for a product's final price by its url_key.
const magentoQuery = `{products(filter:{url_key:{eq:%q}}){items{url_key price_range{minimum_price{final_price{value currency}}}}}}`

// apiURL returns the API request that stands for the product page.
func apiURL(p Product) (string, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return "", err
	}
	switch p.API {
	case WooCommerce:
		api := url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/wp-json/wc/store/v1/products",
			RawQuery: url.Values{"slug": {wooSlug(p.URL)}}.Encode()}
		return api.String(), nil
	case Magento:
		api := url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/graphql",
			RawQuery: url.Values{"query": {fmt.Sprintf(magentoQuery, magentoKey(p.URL))}}.Encode()}
		return api.String(), nil
	}
	return "", fmt.Errorf("unknown API %q", p.API)
}

// fromAPI reads the price from an API response.
func fromAPI(body []byte, p Product) (decimal.Decimal, string, error) {
	switch p.API {
	case WooCommerce:
		return fromWooCommerce(body, wooSlug(p.URL))
	case Magento:
		return fromMagento(body, magentoKey(p.URL))
	}
	return decimal.Zero, "", fmt.Errorf("unknown API %q", p.API)
}

// wooSlug is the product slug in a WooCommerce page URL: its last segment.
func wooSlug(pageURL string) string {
	u, err := url.Parse(pageURL)
	if err != nil {
		return ""
	}
	return path.Base(strings.TrimSuffix(u.Path, "/"))
}

// wooProduct is the part of a WooCommerce Store API product we read.
type wooProduct struct {
	Slug   string `json:"slug"`
	Prices struct {
		Price        string `json:"price"`
		CurrencyCode string `json:"currency_code"`
		MinorUnit    int32  `json:"currency_minor_unit"`
	} `json:"prices"`
}

// fromWooCommerce reads the price of the product with the given slug. Some
// shops ignore the slug filter and list other products, so the slug of the
// answer is checked.
func fromWooCommerce(body []byte, slug string) (decimal.Decimal, string, error) {
	var products []wooProduct
	if err := json.Unmarshal(body, &products); err != nil {
		return decimal.Zero, "", fmt.Errorf("parse: %w", err)
	}
	i := slices.IndexFunc(products, func(p wooProduct) bool { return strings.EqualFold(p.Slug, slug) })
	if i < 0 {
		return decimal.Zero, "", errNoPrice
	}
	pr := products[i].Prices
	minor, err := decimal.NewFromString(pr.Price)
	if err != nil || !minor.IsPositive() {
		return decimal.Zero, "", errNoPrice
	}
	return minor.Shift(-pr.MinorUnit), pr.CurrencyCode, nil
}

// magentoKey is the url_key in a Magento page URL: its last segment, without
// the ".html" suffix.
func magentoKey(pageURL string) string {
	return strings.TrimSuffix(wooSlug(pageURL), ".html")
}

type magentoProduct struct {
	URLKey     string `json:"url_key"`
	PriceRange struct {
		MinimumPrice struct {
			FinalPrice struct {
				Value    json.Number `json:"value"`
				Currency string      `json:"currency"`
			} `json:"final_price"`
		} `json:"minimum_price"`
	} `json:"price_range"`
}

// fromMagento reads the final price of the product with the given url_key.
func fromMagento(body []byte, key string) (decimal.Decimal, string, error) {
	var resp struct {
		Data struct {
			Products *struct {
				Items []magentoProduct `json:"items"`
			} `json:"products"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return decimal.Zero, "", fmt.Errorf("parse: %w", err)
	}
	if len(resp.Errors) > 0 {
		return decimal.Zero, "", fmt.Errorf("graphql: %s", resp.Errors[0].Message)
	}
	if resp.Data.Products == nil {
		return decimal.Zero, "", errNoPrice
	}
	items := resp.Data.Products.Items
	i := slices.IndexFunc(items, func(p magentoProduct) bool { return p.URLKey == key })
	if i < 0 {
		return decimal.Zero, "", errNoPrice
	}
	fp := items[i].PriceRange.MinimumPrice.FinalPrice
	price, err := decimal.NewFromString(fp.Value.String())
	if err != nil || !price.IsPositive() {
		return decimal.Zero, "", errNoPrice
	}
	return price, fp.Currency, nil
}
