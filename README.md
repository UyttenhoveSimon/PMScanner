# PMScanner

Scrapes precious-metal prices from dealer websites and ranks them by price per gram.

```sh
go run .
```

Products are listed in `products.go`. Prices are read from schema.org structured
data (JSON-LD, microdata or `product:price:amount` meta tags), so most shops need
no site-specific code: add a `Product` with its URL and weight. If a shop's first
source holds the wrong price, set `Sources` to pick another (see proaurum).
