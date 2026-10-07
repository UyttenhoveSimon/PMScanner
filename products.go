package main

import (
	"github.com/shopspring/decimal"

	"pmscanner/scraper"
)

var (
	kg   = decimal.NewFromInt(1000)
	oz   = decimal.RequireFromString("31.1035")
	tube = oz.Mul(decimal.NewFromInt(25))
)

// products lists the pages to scan. Any shop page exposing schema.org price
// data (JSON-LD, microdata or Open Graph) works without site-specific code.
//
// Browser: true marks shops behind a bot challenge (goldavenue: Vercel,
// silvergoldbull: Cloudflare); they need a Chromium-based browser installed.
var products = []scraper.Product{
	{Site: "or.fr", Description: "Gold bar 1kg Valcambi", WeightGrams: kg,
		URL: "https://or.fr/achat/or/lingots/lingot-d-or-1-kilogramme-valcambi-60"},
	{Site: "or.fr", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://or.fr/achat/or/pieces/krugerrand-or-1-once-2026-south-african-mint-498"},
	{Site: "or.fr", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://or.fr/achat/or/pieces/maple-leaf-or-1-once-2026-royal-canadian-mint-470"},
	{Site: "or.fr", Description: "Silver Maple Leaf 1oz (tube of 25)", WeightGrams: tube,
		URL: "https://or.fr/achat/argent/pieces/maple-leaf-argent-1-once-2026-tube-25-pieces-royal-canadian-mint-510"},

	{Site: "suissegold", Description: "Gold bar 1kg", WeightGrams: kg,
		URL: "https://suissegold.com/en/product/1-kilogram-gold-bullion-bar-type-of-our-choice"},
	{Site: "suissegold", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://suissegold.com/en/product/1-ounce-oz-2026-south-african-krugerrand"},
	{Site: "suissegold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://suissegold.com/en/product/2026-1-ounce-oz-canadian-maple-leaf-coin"},
	{Site: "suissegold", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://suissegold.com/en/product/2026-1-ounce-oz-canadian-maple-leaf-silver-bullion-coin"},

	{Site: "orobel", Description: "Gold bar 1kg", WeightGrams: kg,
		URL: "https://www.orobel.biz/produit/lingot-1-kg-or"},
	{Site: "orobel", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://www.orobel.biz/produit/krugerrand-1oz"},
	{Site: "orobel", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.orobel.biz/produit/maple-leaf-1-oz"},
	{Site: "orobel", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.orobel.biz/produit/acheter-piece-maple-leaf-argent-en-ligne-orobel"},

	{Site: "degussa", Description: "Gold bar 1kg Degussa", WeightGrams: kg,
		URL: "https://degussa.com/de-de/1000-g-degussa-goldbarren-gegossen-p5376/"},
	{Site: "degussa", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://degussa.com/de-de/1-oz-kruegerrand-goldmuenze-suedafrika-verschiedene-jahrgaenge-p3191/"},
	{Site: "degussa", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://degussa.com/de-de/1-oz-maple-leaf-goldmuenze-50-dollars-kanada-verschiedene-jahrgaenge-p3269/"},
	{Site: "degussa", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://degussa.com/de-de/1-oz-maple-leaf-silbermuenze-5-dollars-kanada-verschiedene-jahrgaenge-p3984/"},

	// proaurum's JSON-LD price is its buy-back price; the selling price is in meta.
	{Site: "proaurum", Description: "Gold bar 1kg Heraeus", WeightGrams: kg, Sources: proaurum,
		URL: "https://www.proaurum.de/shop/goldbarren-1000-gramm-heraeus/"},
	{Site: "proaurum", Description: "Gold Krugerrand 1oz", WeightGrams: oz, Sources: proaurum,
		URL: "https://www.proaurum.de/shop/1-unze-gold-kruegerrand/"},
	{Site: "proaurum", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, Sources: proaurum,
		URL: "https://www.proaurum.de/shop/1-unze-gold-maple-leaf/"},
	{Site: "proaurum", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, Sources: proaurum,
		URL: "https://www.proaurum.de/shop/silbermunze-1-unze-maple-leaf-regelbesteuert-aktueller-jahrgang/"},

	{Site: "achat-or-et-argent", Description: "Gold bar 1kg", WeightGrams: kg,
		URL: "https://www.achat-or-et-argent.fr/or/lingot-1kg/38"},
	{Site: "achat-or-et-argent", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://www.achat-or-et-argent.fr/or/krugerrand/12"},
	{Site: "achat-or-et-argent", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.achat-or-et-argent.fr/or/maple-leaf-1-once-or/3192"},
	{Site: "achat-or-et-argent", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.achat-or-et-argent.fr/argent/maple-leaf-1-once/1668"},

	{Site: "stonexbullion", Description: "Gold bar 1kg", WeightGrams: kg,
		URL: "https://stonexbullion.com/fr/lingots-or/1000g/1-kilo-lingot-d-or-plusieurs-fabricants-lbma/"},
	{Site: "stonexbullion", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://stonexbullion.com/fr/pieces-or/krugerrand/1-oz-krugerrand-or-plusieurs-annees/"},
	{Site: "stonexbullion", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://stonexbullion.com/fr/pieces-or/maple-leaf-du-canada/1-oz-maple-leaf-or-2026/"},
	{Site: "stonexbullion", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://stonexbullion.com/fr/pieces-argent/maple-leaf-du-canada/1-oz-maple-leaf-argent-plusieurs-annees/"},

	{Site: "philoro", Description: "Gold bar 1kg", WeightGrams: kg,
		URL: "https://philoro.de/produkt/goldbarren-1000-g-diverse-hersteller-2711v"},
	{Site: "philoro", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://philoro.de/produkt/gold-kruegerrand-1-oz-diverse-jahrgaenge-1001v"},
	{Site: "philoro", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://philoro.de/produkt/gold-maple-leaf-1-oz-2026-2131v"},
	{Site: "philoro", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://philoro.de/produkt/silber-maple-leaf-1-oz-2026-2136v"},

	{Site: "hollandgold", Description: "Gold bar 1kg Heraeus", WeightGrams: kg,
		URL: "https://www.hollandgold.nl/heraeus-1000-gram-goudbaar.html"},
	{Site: "hollandgold", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://www.hollandgold.nl/krugerrand-1-troy-ounce-gouden-munt.html"},
	{Site: "hollandgold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.hollandgold.nl/maple-leaf-1-troy-ounce-gouden-munt-diverse-jaartallen.html"},
	{Site: "hollandgold", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.hollandgold.nl/maple-leaf-1-troy-ounce-zilver-div-jaartallen.html"},

	{Site: "bdor", Description: "Gold bar 1kg", WeightGrams: kg,
		URL: "https://www.bdor.fr/achat-or-en-ligne/lingot-d-or-1-kg"},
	{Site: "bdor", Description: "Gold Krugerrand 1oz", WeightGrams: oz,
		URL: "https://www.bdor.fr/achat-or-en-ligne/piece-d-or-krugerrand"},
	{Site: "bdor", Description: "Gold Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.bdor.fr/achat-or-en-ligne/piece-d-or-maple-leaf-once"},
	{Site: "bdor", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://www.bdor.fr/achat-or-en-ligne/maple-leaf-argent-1-once"},

	{Site: "goldsilver.be", Description: "Silver Maple Leaf 1oz", WeightGrams: oz,
		URL: "https://goldsilver.be/en/1-oz-30-gr/10863-1-oz-silver-maple-leaf-2026-5-bu.html"},

	{Site: "goldavenue", Description: "Gold bar 1kg PAMP", WeightGrams: kg, Browser: true, Cookies: eur,
		URL: "https://www.goldavenue.com/en/buy/gold/product/1-kg-gold-bar-999-9-fine-gold-carbon-measured-pamp-suisse"},
	{Site: "goldavenue", Description: "Gold Krugerrand 1oz", WeightGrams: oz, Browser: true, Cookies: eur,
		URL: "https://www.goldavenue.com/en/buy/gold/product/1-oz-fine-gold-coin-916-7-krugerrand-mixed-years"},
	{Site: "goldavenue", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, Browser: true, Cookies: eur,
		URL: "https://www.goldavenue.com/en/buy/gold/product/1-oz-fine-gold-coin-999-9-maple-leaf-bu-mixed-years"},
	{Site: "goldavenue", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, Browser: true, Cookies: eur,
		URL: "https://www.goldavenue.com/en/buy/silver/product/1-ounce-silver-coin-canada-maple-leaf-2026"},

	// silvergoldbull has no structured price data; read the "1+" row of the
	// quantity price table, or the headline price when there is no table.
	{Site: "silvergoldbull", Description: "Gold bar 1kg RCM", WeightGrams: kg, Browser: true, Selectors: sgbPrice,
		URL: "https://silvergoldbull.be/1-kilo-gold-bar-royal-canadian-mint"},
	{Site: "silvergoldbull", Description: "Gold Krugerrand 1oz", WeightGrams: oz, Browser: true, Selectors: sgbPrice,
		URL: "https://silvergoldbull.be/1-oz-2026-krugerrand-gold-coin-rand-refinery"},
	{Site: "silvergoldbull", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, Browser: true, Selectors: sgbPrice,
		URL: "https://silvergoldbull.be/1-oz-2026-canadian-maple-leaf-gold-coin-royal-canadian-mint"},
	{Site: "silvergoldbull", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, Browser: true, Selectors: sgbPrice,
		URL: "https://silvergoldbull.be/1-oz-2026-canadian-maple-leaf-silver-coin-royal-canadian-mint"},
}

var (
	proaurum = []scraper.Source{scraper.Meta}
	// goldavenue defaults to CHF.
	eur      = map[string]string{"currency": "EUR"}
	sgbPrice = []string{
		`div[class~="tw:grid-cols-10"] > div:containsOwn("1+") + div`,
		`h1 + div > div[class~="tw:font-bold"]`,
	}
)
