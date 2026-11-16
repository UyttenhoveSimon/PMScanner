# Maintenance notes

What we learned about individual shops, how prices were verified, and what is
still open. Last updated 2026-11-16.

## Bot walls and blocks

- **gvs group** (goldvorsorge.at/.ch, zlatosrebro.hr, aranykereskedes, gvs.ee,
  kupnozlota): since mid-October 2026 they answer with a proof-of-work
  "Sicherheitsprüfung" page (HTTP 200, no price → "no price found"). Now
  `Browser: true` with `RateGroup: "gvs"`. After a test scan with 8 parallel
  tabs, every request from our IP got HTTP 403 (likely an IP block); since then
  browser tabs also respect the 2-per-group limit. **To do:** rescan once the
  block lifts and check the ~31 pages that timed out.
- **esg**: HTTP 403 to plain HTTP; `Browser: true`. Only JSON-LD, no API.
- **centerzlata**: HTTP 403 to plain HTTP; read through the WooCommerce Store
  API instead (no browser needed).

## JSON APIs

Survey of November 2026 (details in the README, "For shop owners"):

- **WooCommerce Store API** `/wp-json/wc/store/v1/products?slug=…`: ~31 shops.
  Products switched only where the API price equals the page price (207
  products, 28 shops). Prices are in minor units (`currency_minor_unit`).
  - comprarlingotes ignores `slug` and returns the same product every time —
    `fromWooCommerce` checks the slug; not switched.
  - vitusguld: API returns nothing.
  - Missing in the API (kept on page): avangardgold Silver Maple and Vreneli,
    italpreziosi Maple Leaf and Vreneli, norskmynthandel
    Philharmonic/Britannia/Kangaroo, orodainvestimento Eagle.
  - API differs from page (kept on page): orohispanica silver bar.
  - centerzlata's Store API excludes VAT (1748.91 × 1.22 = 2133.67), so its
    silver bar is `VATExcluded`.
- **Magento GraphQL** `/graphql`: argentor and bessergold work and are switched.
  reisebank: url_key doesn't match our URLs. smh: "index does not exist yet".
  orocash, ourocaixa, tsbnk, hattongardenmetals: false positives (HTML).

## VAT verification (November 2026)

Only `VATExcluded` changes the number (VAT is added in `tax.go`); the other
labels only change what the page shows.

Fixed:

- zlatopoposti, svetplemenitihkovin Silver Maple Leaf: JSON-LD was ex-VAT, page
  shows incl. 22% DDV → read through the Store API (69.44 / 62.04).
- comprarlingotes Silver Maple Leaf: page says "(+IVA)" → `VATExcluded`.
- **suissegold**: price and currency depend on the visitor's location. From
  Europe: CHF incl. 8.1% VAT. From the US (GitHub Actions): USD without VAT.
  Proof: CHF/USD = 0.8389 for gold (exchange rate, gold is VAT exempt) and
  0.9069 = 0.8389 × 1.081 for other metals. Fix: `ExportCurrency: "USD"` —
  prices in that currency count as VAT excluded. Country/currency is stored in
  the session (form POST); setting it by POST didn't work. Older saved scans
  keep the too-low prices.

Checked, correct (5 cheapest VAT-included per non-gold category):

- Page text says incl. VAT: philoro (.li/.ch/.at), geiger, geiger.ch, degussa,
  proaurum, auragentum, anlagegold24, eurodata.lt. Swiss ones look cheap
  because Swiss VAT is 8.1%.
- baird: consistent with UK VAT by spot math. tavex.pl: silver spot 195.35
  PLN/oz + 23% VAT + 2.6% premium. mennicaskarbowa: assumed gross.
- 37 non-gold pages flag `valueAddedTaxIncluded: false` in JSON-LD: 22 were
  already excluded/margin; the 12 marked included (orohispanica, bleyer, aurom,
  neogold, nyfortuna, tavast, topgold, zlatypristav) are priced at or above the
  category median, so judged fine, but not verified on page text.

Label only (price is right, label could be better):

- hollandgold, silvermountain Silver Maple: margin scheme ("Margeregeling").
- goudwisselkantoor Palladium coin: margin scheme.
- norskmynthandel silver: "Fritatt MVA" (exempt).
- goldsilver.be: terms say all taxes included; €/g suggests margin scheme.

Uncertain:

- smart-elements (Rhodium, Osmium): scan 8392/925 "USD" vs API 8030/885 "CHF"
  (ratio 0.957); neither states VAT. Ask the shop.

## Other open errors (scan 2026-11-15)

- zoloto-md: 12 × "no price found" since ~2026-11-05; not investigated.
- HTTP 404, URLs need updating: philoro.ch, norskmynthandel, vitusguld,
  auracasa, atkinsons, philoro, mygold.bg.
- xrisos.gr: HTTP 429. azergold: TLS certificate error.
- ~31 flaky errors that come and go.
