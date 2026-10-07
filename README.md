# PMScanner

Scrapes precious-metal prices from dealer websites and ranks them by price per gram.

```sh
go run .               # print a table
go run . -serve :8080  # web page on http://localhost:8080, rescans every 30 min
```

The web page also serves the latest results as JSON at `/api/prices`.
Requires Go 1.27.

Products are listed in `products.go`. Prices are read from schema.org structured
data (JSON-LD, microdata or `product:price:amount` meta tags), so most shops need
no site-specific code: add a `Product` with its URL and weight. If a shop's first
source holds the wrong price, set `Sources` to pick another (see proaurum).

Shops without structured data can use `Selectors` (CSS selectors for the price
element). Shops behind a bot challenge set `Browser: true` and are loaded in a
headless Chromium-based browser (Chrome, Chromium, Brave or Edge, found
automatically; override with `PMSCANNER_BROWSER=/path/to/browser`).
