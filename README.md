# PMScanner

[![CI](https://github.com/UyttenhoveSimon/PMScanner/actions/workflows/ci.yml/badge.svg)](https://github.com/UyttenhoveSimon/PMScanner/actions/workflows/ci.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/UyttenhoveSimon/PMScanner)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Compare precious-metal prices across European bullion dealers, ranked by price per gram.**

PMScanner scrapes gold and silver product pages from many dealers at once and tells
you who sells cheapest, either as a terminal table or as a small self-hosted web page.

```
SITE           PRODUCT                PRICE          PER GRAM  URL
or.fr          Gold bar 1kg Valcambi  119856.23 EUR  119.86    https://or.fr/achat/or/lingots/...
hollandgold    Gold bar 1kg Heraeus   119989.86 EUR  119.99    https://www.hollandgold.nl/...
stonexbullion  Gold Krugerrand 1oz    3730.60 EUR    119.94    https://stonexbullion.com/fr/...
...
```

## Features

- **No per-site scraping code.** Prices come from the schema.org data shops publish
  for search engines (JSON-LD, microdata, Open Graph), so redesigns rarely break it.
- **Fast.** All pages are fetched concurrently; a full scan takes about 25 seconds.
- **Works behind bot walls.** Shops protected by Cloudflare or similar are loaded in
  a headless Chromium-based browser.
- **Web page and JSON API.** `-serve` starts a page that rescans in the background
  and exposes the results at `/api/prices`.
- **Easy to extend.** Adding a shop is usually one line in [`products.go`](products.go).

## Quick start

Requires [Go 1.27+](https://go.dev/dl/).

```sh
git clone https://github.com/UyttenhoveSimon/PMScanner.git
cd PMScanner
go run .
```

### Web page

```sh
go run . -serve :8080
```

Open <http://localhost:8080>. Prices are grouped by product with the cheapest
offer highlighted, and rescanned every 30 minutes.

| Flag       | Default | Description                                          |
| ---------- | ------- | ---------------------------------------------------- |
| `-serve`   |         | Address to serve the web page on, e.g. `:8080`       |
| `-refresh` | `30m`   | How often the web page rescans prices                |

| Endpoint      | Description                     |
| ------------- | ------------------------------- |
| `/`           | Price tables                    |
| `/api/prices` | Latest scan as JSON             |

### Browser mode

Products marked `Browser: true` need a Chromium-based browser (Chrome, Chromium,
Brave or Edge). It is found automatically; to use another one:

```sh
PMSCANNER_BROWSER="/path/to/browser" go run .
```

## Adding a shop

Add an entry to [`products.go`](products.go):

```go
{Site: "myshop", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
    URL: "https://myshop.example/krugerrand-1oz"},
```

Start the description with one of the product categories (`Gold bar 1kg`,
`Gold Krugerrand 1oz`, `Gold Maple Leaf 1oz`, `Silver Maple Leaf 1oz`) to group it
on the web page. If the price doesn't come out right, these options help:

| Field       | Use when                                                                  |
| ----------- | ------------------------------------------------------------------------- |
| `Sources`   | The shop's first data source holds the wrong price (e.g. a buy-back price) |
| `Selectors` | The shop has no structured data; give CSS selectors for the price element  |
| `Cookies`   | The shop needs a cookie, e.g. to show prices in EUR                        |
| `Browser`   | The shop blocks plain HTTP clients with a JavaScript challenge            |

## How it works

For each product page, PMScanner tries these sources in order and keeps the first price found:

1. `Selectors`, if the product defines any
2. JSON-LD: `<script type="application/ld+json">` with a schema.org `Offer`
3. Microdata: `itemprop="price"`
4. Open Graph: `<meta property="product:price:amount">`

```
.
├── main.go          # CLI: scan and print the table
├── web.go, web.html # -serve: web page and JSON API
├── products.go      # the list of shops and products
└── scraper/         # fetching (HTTP or headless browser) and price extraction
```

## Contributing

Issues and pull requests are welcome, especially new dealers. Before opening a PR:

```sh
gofmt -l . && go vet ./... && go test ./...
```

## Disclaimer

Prices are scraped from public pages and may be delayed or wrong. Always check the
dealer's site before buying. Respect each site's terms of use.

## License

[MIT](LICENSE)
