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
	// Belgium
	{Site: "goudwisselkantoor", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.shop-goudwisselkantoor.be/goud-kopen/staven/goudstaaf-1000-gram-diverse-producenten"},
	{Site: "goudwisselkantoor", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.shop-goudwisselkantoor.be/goud-kopen/munten/gouden-krugerrand-1-oz-divers-jaar"},
	{Site: "goudwisselkantoor", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.shop-goudwisselkantoor.be/goud-kopen/munten/gouden-maple-leaf-1-oz-divers-jaar"},
	{Site: "goudwisselkantoor", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.shop-goudwisselkantoor.be/zilver-kopen/munten/zilveren-maple-leaf-1-oz-divers-jaar"},

	// Netherlands
	{Site: "goudonline", Description: "Gold bar 1kg C.Hafner", WeightGrams: kg, URL: "https://goudonline.nl/goud-kopen/c-hafner-goudbaar-1000-gram/"},
	{Site: "goudonline", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://goudonline.nl/gouden-munten-kopen/1-troy-ounce-gouden-krugerrand/"},
	{Site: "goudonline", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://goudonline.nl/goud-kopen/1-troy-ounce-gouden-maple-leaf-kopen/"},
	{Site: "goudonline", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://goudonline.nl/zilveren-munten/1-troy-ounce-zilveren-maple-leaf-munt/"},
	{Site: "silvermountain", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.thesilvermountain.nl/en/1-kilo-gold-bar"},
	{Site: "silvermountain", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.thesilvermountain.nl/en/1-troy-ounce-gold-krugerrand-coin"},
	{Site: "silvermountain", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.thesilvermountain.nl/en/1-troy-ounce-gold-maple-leaf"},
	{Site: "silvermountain", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.thesilvermountain.nl/en/1-troy-ounce-silver-coin-maple-leaf"},

	// France
	{Site: "oretchange", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.oretchange.com/lingots-or/183-achat-lingot-1-kilo-d-or.html"},
	{Site: "oretchange", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.oretchange.com/pieces-or/186-achat-krugerrand-afrique-du-sud.html"},
	{Site: "oretchange", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.oretchange.com/acheter-onces-or/235-682195-maple-leaf-1-once-or.html"},
	{Site: "oretchange", Description: "Silver Maple Leaf 1oz (min. 5)", WeightGrams: oz, URL: "https://www.oretchange.com/pieces-d-argent-internationales/208-achat-maple-leaf-1-once-argent.html"},
	{Site: "or-investissement", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://or-investissement.fr/achat-lingot-or-investissement/1041-achat-lingot-or-1-kilo.html"},
	{Site: "or-investissement", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://or-investissement.fr/achat-piece-or-investissement/16-achat-piece-or-krugerrand.html"},

	// United Kingdom (GBP)
	{Site: "atkinsons", Description: "Gold bar 1kg Metalor", WeightGrams: kg, URL: "https://atkinsonsbullion.com/gold/gold-bars/1kg-gold-bars/metalor-1kg-gold-cast-bar"},
	{Site: "atkinsons", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://atkinsonsbullion.com/gold/gold-coins/1oz-gold-coins/2026-south-african-krugerrand-1oz-gold-coin"},
	{Site: "atkinsons", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://atkinsonsbullion.com/gold/gold-coins/1oz-gold-coins/2026-canadian-maple-leaf-1oz-gold-coin"},
	{Site: "bleyer", Description: "Gold bar 1kg Baird", WeightGrams: kg, URL: "https://www.bleyerbullion.co.uk/shop/gold/gold-bars/1kg-gold-baird-co-cast-bar/"},
	{Site: "bleyer", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.bleyerbullion.co.uk/shop/gold/gold-coins/canadian-gold-coins/1oz-gold-canadian-maple-leaf-coin-2026/"},
	{Site: "bleyer", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.bleyerbullion.co.uk/shop/silver/silver-coins/canadian-silver-coins/1oz-silver-canadian-maple-leaf-coin-2026/"},

	// Germany
	{Site: "kettner", Description: "Gold bar 1kg Umicore", WeightGrams: kg, URL: "https://www.kettner-edelmetalle.de/goldbarren/1-kg-Goldbarren/1011824-1kg-goldbarren-umicore-gussbarren.html"},
	{Site: "kettner", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.kettner-edelmetalle.de/Goldmuenzen/Kruegerrand/21-1-unze-gold-kruegerrand.html"},
	{Site: "kettner", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.kettner-edelmetalle.de/Goldmuenzen/Maple-Leaf/1030403-1-unze-gold-maple-leaf-2026.html"},
	{Site: "kettner", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.kettner-edelmetalle.de/Silbermuenzen/Maple-Leaf/1030340-1-unze-silber-maple-leaf-2026.html"},
	{Site: "anlagegold24", Description: "Gold bar 1kg C.Hafner", WeightGrams: kg, URL: "https://www.anlagegold24.de/1000-Gramm-Goldbarren-CHafner.html"},
	{Site: "anlagegold24", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.anlagegold24.de/1-oz-Gold-Kruegerrand-verschiedene-Jahrgaenge-Deutsch.html"},
	{Site: "anlagegold24", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.anlagegold24.de/1-oz-Gold-Maple-Leaf-Verschiedene-Jahrgaenge.html"},
	{Site: "anlagegold24", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.anlagegold24.de/1-oz-Silber-Maple-Leaf-verschiedene-Jahrgaenge.html"},
	{Site: "esg", Description: "Gold bar 1kg Heraeus", WeightGrams: kg, URL: "https://www.edelmetall-handel.de/goldbarren-1000g-heraeus-gegossen-01011102"},
	{Site: "esg", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.edelmetall-handel.de/goldmuenze-1oz-kruegerrand-aktueller-jahrgang-suedafrika-10000011"},
	{Site: "esg", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.edelmetall-handel.de/goldmuenze-1oz-maple-leaf-aktueller-jahrgang-kanada-10000019"},
	{Site: "esg", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.edelmetall-handel.de/silbermuenze-1oz-maple-leaf-aktueller-jahrgang-19-kanada-10208301"},
	{Site: "reisebank", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://edelmetalle.reisebank.de/products/1000-g-goldbarren/1000ggoldbarren"},
	{Site: "reisebank", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://edelmetalle.reisebank.de/products/1-oz-kruegerrand-goldmuenze/1ozkruegerrandgoldmuenze"},
	{Site: "reisebank", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://edelmetalle.reisebank.de/products/1-oz-maple-leaf-goldmuenze/1ozmapleleafgoldmuenze"},
	{Site: "reisebank", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://edelmetalle.reisebank.de/products/1-oz-maple-leaf-silbermuenze/1ozmapleleafsilbermuenze"},
	{Site: "auragentum", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://auragentum.de/products/1-000-gramm-goldbarren"},
	{Site: "auragentum", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://auragentum.de/products/1-unze-goldmuenze-kruegerrand"},
	{Site: "auragentum", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://auragentum.de/products/1-unze-goldmuenze-maple-leaf"},
	{Site: "auragentum", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://auragentum.de/products/1-unze-silbermuenze-maple-leaf"},
	{Site: "geiger", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.geiger-edelmetalle.de/kruegerrand-2026-suedafrika-1-oz-goldmuenze-1211027.html"},
	{Site: "geiger", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.geiger-edelmetalle.de/maple-leaf-2026-kanada-1-oz-goldmuenze-1211017.html"},
	{Site: "geiger", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.geiger-edelmetalle.de/maple-leaf-2026-kanada-1-oz-silbermuenze-1110114.html"},
	{Site: "apollo", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://apollo-edelmetalle.de/index.php?id_product=54&rewrite=1-x-1-oz-gold-kruegerrand&controller=product"},

	// Austria
	{Site: "philoro.at", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://philoro.at/produkt/goldbarren-1000-g-diverse-hersteller-2711v"},
	{Site: "philoro.at", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://philoro.at/produkt/gold-kruegerrand-11-oz-1001v"},
	{Site: "philoro.at", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://philoro.at/produkt/gold-maple-leaf-11-oz-1010v"},
	{Site: "philoro.at", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://philoro.at/produkt/silber-maple-leaf-1-oz-3022v"},
	{Site: "goldvorsorge.at", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.goldvorsorge.at/1kg-goldbarren-argor-heraeus.html"},
	{Site: "goldvorsorge.at", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.goldvorsorge.at/1-oz-gold-krugerrand-diverse.html"},
	{Site: "goldvorsorge.at", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.goldvorsorge.at/1-oz-gold-maple-leaf-diverse.html"},
	{Site: "goldvorsorge.at", Description: "Silver Maple Leaf 1oz (circulated)", WeightGrams: oz, URL: "https://www.goldvorsorge.at/1-oz-silbermuenze-maple-leaf-zirkuliert.html"},
	{Site: "goldundco", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.goldundco.at/shop/goldmuenze-kruegerrand-1-unze/"},
	{Site: "goldundco", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.goldundco.at/shop/goldmuenze-maple-leaf-1-unze/"},
	{Site: "goldundco", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.goldundco.at/shop/silbermuenze-maple-leaf-1-unze-differenzbesteuert/"},
	{Site: "smh", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.smh.net/krugerrand-1-unze-gold/"},
	{Site: "smh", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.smh.net/maple-leaf-1-unze-gold/"},
	{Site: "smh", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.smh.net/maple-leaf-1-unze-silber/"},
	{Site: "oegussa", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.oegussa.at/de/shop/produktdetail/feingold-barren-1000g"},
	{Site: "muenzeoesterreich", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.muenzeoesterreich.at/produkte/1-kilo-goldbarren-muenze-oesterreich-gussbarren"},

	// Switzerland and Liechtenstein (CHF)
	{Site: "philoro.ch", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://philoro.ch/produkt/goldbarren-1000-g-diverse-hersteller-2711v"},
	{Site: "philoro.ch", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://philoro.ch/produkt/gold-kruegerrand-1-oz-1001v"},
	{Site: "philoro.ch", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://philoro.ch/produkt/gold-maple-leaf-1-oz-1010v"},
	{Site: "philoro.ch", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://philoro.ch/produkt/silber-maple-leaf-1-oz-3022v"},
	// proaurum.ch's JSON-LD and meta prices are buy-back prices.
	{Site: "proaurum.ch", Description: "Gold bar 1kg", WeightGrams: kg, Sources: cssOnly, Selectors: proaurumCH, URL: "https://proaurum.ch/1000-gramm-goldbarren/"},
	{Site: "proaurum.ch", Description: "Gold Krugerrand 1oz", WeightGrams: oz, Sources: cssOnly, Selectors: proaurumCH, URL: "https://proaurum.ch/1-unze-gold-kruegerrand/"},
	{Site: "proaurum.ch", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, Sources: cssOnly, Selectors: proaurumCH, URL: "https://proaurum.ch/1-unze-gold-maple-leaf/"},
	{Site: "proaurum.ch", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, Sources: cssOnly, Selectors: proaurumCH, URL: "https://proaurum.ch/1-unze-silber-maple-leaf/"},
	// goldvorsorge.ch's JSON-LD price is wrong; meta matches the page.
	{Site: "goldvorsorge.ch", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, Sources: metaOnly, URL: "https://goldvorsorge.ch/1kg-goldbarren-argor-heraeus.html"},
	{Site: "goldvorsorge.ch", Description: "Gold Krugerrand 1oz", WeightGrams: oz, Sources: metaOnly, URL: "https://goldvorsorge.ch/1-oz-gold-krugerrand-diverse.html"},
	{Site: "goldvorsorge.ch", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, Sources: metaOnly, URL: "https://goldvorsorge.ch/1-oz-gold-maple-leaf-diverse.html"},
	{Site: "goldvorsorge.ch", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, Sources: metaOnly, URL: "https://goldvorsorge.ch/1-oz-silbermuenze-maple-leaf.html"},
	{Site: "geiger.ch", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.geiger-edelmetalle.ch/goldbarren-1-kg-.9999-argor-heraeus-kein-postversand-1221060.html"},
	{Site: "geiger.ch", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.geiger-edelmetalle.ch/kruegerrand-suedafrika-1-oz-goldmuenze-1211016.html"},
	{Site: "geiger.ch", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.geiger-edelmetalle.ch/maple-leaf-kanada-1-oz-goldmuenze-1211017.html"},
	{Site: "geiger.ch", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.geiger-edelmetalle.ch/maple-leaf-kanada-1-oz-silbermuenze-aktuelle-neuware-sofort-verfuegbar-1111000.html"},
	{Site: "philoro.li", Description: "Gold bar 1kg philoro", WeightGrams: kg, URL: "https://philoro.li/produkt/goldbarren-1000-g-philoro-2501v"},
	{Site: "philoro.li", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://philoro.li/produkt/gold-kruegerrand-1-oz-1001v"},
	{Site: "philoro.li", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://philoro.li/produkt/gold-maple-leaf-1-oz-1010v"},
	{Site: "philoro.li", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://philoro.li/produkt/silber-maple-leaf-1-oz-3022v"},

	// Finland
	{Site: "tavex.fi", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, URL: "https://tavex.fi/kulta/1000g-gold-bar/"},
	{Site: "tavex.fi", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.fi/kulta/1oz-south-african-krugerrand/"},
	{Site: "tavex.fi", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.fi/kulta/1oz-canadian-maple-leaf-gold-coin/"},
	{Site: "tavex.fi", Description: "Silver Maple Leaf 1oz (circulated)", WeightGrams: oz, URL: "https://tavex.fi/hopea/1-oz-kanadan-maple-leaf-hopeakolikot/"},
	{Site: "jalonom", Description: "Gold bar 1kg Umicore", WeightGrams: kg, URL: "https://www.jalonom.com/en/shop/gold-bar-1000g-umicore-576"},
	{Site: "jalonom", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.jalonom.com/en/shop/gold-coins-6/gold-coin-1oz-krugerrand-south-africa-606"},
	{Site: "jalonom", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.jalonom.com/en/shop/silver-coins-7/silver-coin-1oz-maple-leaf-canada-628"},

	// Sweden (SEK)
	{Site: "tavex.se", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://tavex.se/guld/1-kg-guldtacka-olika-tillverkare/"},
	{Site: "tavex.se", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.se/guld/1oz-sydafrikansk-krugerrand/"},
	{Site: "tavex.se", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.se/guld/1oz-kanadensisk-maple-leaf-guldmynt/"},
	{Site: "tavex.se", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.se/silver/1oz-kanadensiskt-maple-leaf-silvermynt/"},

	// Norway (NOK)
	{Site: "tavex.no", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.no/gull/1oz-south-african-krugerrand/"},
	{Site: "tavex.no", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.no/gull/kanadisk-maple/"},
	{Site: "tavex.no", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.no/solv/1-oz-kanadiske-maple-leaf-solvmynt/"},
	{Site: "norskmynthandel", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://norskmynthandel.no/product/2026-sor-afrika-1-oz-gullmynt-krugerrand-bu-m-air-tite-kapsel/"},
	{Site: "norskmynthandel", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://norskmynthandel.no/product/2026-canada-1-oz-gullmynt-maple-leaf-bu-m-air-tite-kapsel/"},
	{Site: "norskmynthandel", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://norskmynthandel.no/product/2026-canada-1-oz-solvmynt-maple-leaf-bu/"},

	// Denmark (DKK)
	{Site: "tavex.dk", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://tavex.dk/guld/1-kg-guldbarre-forskellige-maerker/"},
	{Site: "tavex.dk", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.dk/guld/1oz-sydafrikansk-krugerrand/"},
	{Site: "tavex.dk", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.dk/guld/1oz-canadisk-maple-leaf-guldmont/"},
	{Site: "tavex.dk", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.dk/solv/1oz-canadisk-maple-leaf-soelvmoent/"},
	{Site: "vitusguld", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://vitusguld.dk/produkt/1000-gr-stoebt-guldbarre-9999-%e2%80%b0-argor-heraeus-schweiz/"},
	{Site: "vitusguld", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://vitusguld.dk/produkt/1-oz-sydafrikansk-krugerrand-guldmoent-311-gr-24-karat-9999-%e2%80%b0-aar-2026/"},
	{Site: "vitusguld", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://vitusguld.dk/produkt/1-oz-canadian-maple-leaf-guldmoent-311-gr-24-karat-9999-%e2%80%b0-aar-2026/"},
	{Site: "vitusguld", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://vitusguld.dk/produkt/1-oz-canadian-maple-leaf-soelvmoent-311-gr-finsoelv-9999-%e2%80%b0-aar-2025/"},
	{Site: "nyfortuna", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://nyfortuna.dk/produkt/guldbarre-stoebt-1000g/"},
	{Site: "nyfortuna", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://nyfortuna.dk/produkt/krugerrand-1-oz/"},
	{Site: "nyfortuna", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://nyfortuna.dk/produkt/maple-leaf-1oz/"},
	// Spain
	{Site: "andorrano", Description: "Gold bar 1kg C.Hafner", WeightGrams: kg, URL: "https://www.andorrano-joyeria.com/lingotes-de-oro/lingote-oro-c-hafner-1kg-info"},
	{Site: "andorrano", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.andorrano-joyeria.com/tienda/monedas-de-oro/krugerrand-2026-1oz-oro-info"},
	{Site: "andorrano", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.andorrano-joyeria.com/tienda/monedas-de-oro/maple-leaf/maple-leaf-2026-1oz-oro-info"},
	{Site: "andorrano", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.andorrano-joyeria.com/tienda/monedas-de-plata/canada/maple-leaf-2026-1oz-plata-info"},
	{Site: "orodeinversion", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://orodeinversion.com/shop/lingote-de-oro-1-kg-142"},
	{Site: "orodeinversion", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://orodeinversion.com/shop/krugerrand-oro-1-oz-20"},
	{Site: "orodeinversion", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://orodeinversion.com/shop/maple-leaf-oro-1-oz-19"},
	{Site: "orodeinversion", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://orodeinversion.com/shop/maple-leaf-plata-1-oz-22"},
	{Site: "dracma", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.dracmametales.com/monedas-de-oro/1oz-canada-maple-leaf-2026-bu-oro"},
	{Site: "ladobla", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://ladoblabullion.com/producto/lingote-de-oro-1000-g-argor-heraeus/"},
	{Site: "ladobla", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://ladoblabullion.com/producto/1-onza-oro-krugerrand-sudafrica/"},
	{Site: "orohispanica", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.orohispanica.es/tienda/lingotes-de-oro/lingote-de-oro-fino-9999-de-1000-g-en-blister-argor-heraeus/"},
	{Site: "orohispanica", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.orohispanica.es/tienda/monedas-de-oro/moneda-krugerrand-sudafricano-1-oz-oro-puro/"},
	{Site: "orohispanica", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.orohispanica.es/tienda/monedas-de-oro/moneda-maple-leaf-1-oz-oro-puro/"},
	{Site: "comprarlingotes", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.comprarlingotes.com/tienda/oro/lingotes-de-oro/lingote-de-oro-de-1000-gramos-1-kg-de-argor-heraeus/"},
	{Site: "comprarlingotes", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.comprarlingotes.com/tienda/oro/monedas-de-oro/moneda-maple-leaf-canadiense-de-oro-1-oz/"},
	{Site: "comprarlingotes", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.comprarlingotes.com/tienda/plata/monedas-de-plata/moneda-maple-leaf-de-1-onza-de-plata-2019/"},
	{Site: "ciode", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://ciode.es/producto/lingote-de-oro-1000-grs/"},
	{Site: "ciode", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://ciode.es/producto/krugerrand-sudafrica-1-oz-oro-varios-a-os/"},
	{Site: "ciode", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://ciode.es/producto/maple-oro-50-canada-31-10-grs-varios-a-os-2/"},
	{Site: "ciode", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://ciode.es/producto/onza-plata-maple-leaf-2017-canada/"},

	// Italy
	{Site: "italpreziosi", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://shop.italpreziosi.it/investimento/lingotti/oro/lingotto-oro-1-kg/"},
	{Site: "italpreziosi", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://shop.italpreziosi.it/investimento/monete/oro/1-oz-maple-leaf-oro-2026/"},
	{Site: "orodainvestimento", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.orodainvestimento.it/acquisto/lingotto-doro-1-kg-lbma/"},
	{Site: "orodainvestimento", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.orodainvestimento.it/acquisto/krugerrand-oro-fdc-2026-1-oz/"},
	{Site: "orodainvestimento", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.orodainvestimento.it/acquisto/maple-leaf-doro-2026-1-oz/"},
	{Site: "orodainvestimento", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.orodainvestimento.it/acquisto/maple-leaf-argento-2026-fdc-1-oz/"},
	{Site: "euronummus", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.euronummus.it/45-sudafrica-krugerrand-oro-oncia.html"},
	{Site: "euronummus", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.euronummus.it/210-canada-foglia-acero-50-dollari.html"},
	{Site: "investioro", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.investioro.it/prodotto/lingotto-oro-puro-1-kilo/"},
	{Site: "investioro", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.investioro.it/prodotto/krugerrand-oro-1-oncia-31-1/"},
	{Site: "orocash", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://invest.orocash.it/product/Lingotto-da-fusione-1000-grammi/226338/"},

	// Portugal
	{Site: "ourocaixa", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://invest.ourocaixa.pt/product/Lingotto-da-fusione-1000-grammi/334078/"},

	// Greece
	{Site: "xrisos.gr", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://xrisos.gr/products/maple-leaf-silver-coin-1-oz"},

	// Cyprus
	{Site: "kdggold", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.kdggold.com/en/gold-coins/krugerrand-1oz-gold-mixed-years-en/"},
	{Site: "kdggold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.kdggold.com/en/gold-coins/canadian-maple-leaf-1-oz-2026-gold/"},

	// Slovenia
	{Site: "centerzlata", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://centerzlata.com/izdelek/argor-heraeus-zlata-palica-od-1-kg/"},
	{Site: "centerzlata", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://centerzlata.com/izdelek/krugerrand-1-unca-zlati-kovanec/"},
	{Site: "centerzlata", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://centerzlata.com/izdelek/gold-maple-leaf-1-unca/"},
	{Site: "goldstore.si", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://goldstore.si/prodaja/zlate-palice/1000-g-zlate-palice"},
	{Site: "goldstore.si", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://goldstore.si/prodaja/zlati-kovanci/juzna-afrika/krugerrand-1-oz"},
	{Site: "goldstore.si", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://goldstore.si/prodaja/zlati-kovanci/kanada/kanadski-javorjev-list-1oz"},
	{Site: "hisazlata", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://hisazlata.com/izdelek/zlata-palica-1000-g-argor-heraeus/"},
	{Site: "zlatopoposti", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://zlatopoposti.si/izdelek/1000-g-zlata-palica/"},
	{Site: "zlatopoposti", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://zlatopoposti.si/izdelek/krugerrand-zlatnik-1-oz/"},
	{Site: "zlatopoposti", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://zlatopoposti.si/izdelek/javorjev-list-zlati-1-oz/"},
	{Site: "zlatopoposti", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://zlatopoposti.si/izdelek/javorjev-list-srebrni-1-unca/"},
	{Site: "svetplemenitihkovin", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.svetplemenitihkovin.si/shop/zlate-palice/zlata-palica-1000g/"},
	{Site: "svetplemenitihkovin", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.svetplemenitihkovin.si/shop/afriski-krugerrand/33931-g-zlati-juznoafriski-krugerrand/"},
	{Site: "svetplemenitihkovin", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.svetplemenitihkovin.si/shop/kanadski-list1/3110g-srebrni-kanadski-javorjev-list/"},

	// Croatia
	{Site: "zlatosrebro.hr", Description: "Gold bar 1kg Heraeus (bonded warehouse)", WeightGrams: kg, URL: "https://www.zlatosrebro.hr/1-kg-zlatna-poluga-heraeus-u-bescarinskom-skladistu.html"},
	{Site: "zlatosrebro.hr", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.zlatosrebro.hr/1-unca-zlatnik-krugerrand-razni.html"},
	{Site: "zlatosrebro.hr", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.zlatosrebro.hr/1-unca-zlatnik-javorov-list-razni.html"},
	{Site: "zlatosrebro.hr", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.zlatosrebro.hr/1-unca-srebrnjak-javorov-list-2022.html"},

	// Poland
	{Site: "mennicaskarbowa", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.mennicaskarbowa.pl/zlote-sztabki/632-1000-g-1-kg-sztabka-zlota-wysylka-24-h.html"},
	{Site: "mennicaskarbowa", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.mennicaskarbowa.pl/zlote-monety/1348-zlota-moneta-krugerrand-1-uncja-nowe-roczniki-20252026-wysylka-24-h.html"},
	{Site: "mennicaskarbowa", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.mennicaskarbowa.pl/zlote-monety/1345-zlota-moneta-kanadyjski-lisc-klonowy-1-uncja-nowe-roczniki-20252026-wysylka-24-h.html"},
	{Site: "mennicaskarbowa", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.mennicaskarbowa.pl/srebrne-monety/216-srebrna-moneta-kanadyjski-lisc-klonowy-1-uncja-wysylka-24-h.html"},
	{Site: "tavex.pl", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://tavex.pl/zloto/zlota-sztabka-1-kg/"},
	{Site: "tavex.pl", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.pl/zloto/zloty-krugerrand-1-oz/"},
	{Site: "tavex.pl", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.pl/zloto/1oz-kanadyjski-lisc-klonu-zlota-moneta/"},
	{Site: "tavex.pl", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.pl/srebro/srebrny-kanadyjski-lisc-klonu-1-oz/"},
	{Site: "mennicainwestorow", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://mennicainwestorow.pl/sztabka-zlota-1kg"},
	{Site: "mennicainwestorow", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://mennicainwestorow.pl/krugerrand-1-uncja-zlota-wysylka-24h"},
	{Site: "mennicainwestorow", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://mennicainwestorow.pl/kanadyjski-lisc-klonowy-1-uncja-zlota"},
	{Site: "metalmarket", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.metalmarket.eu/en/products/gold-bar-1000-grams-23"},
	{Site: "metalmarket", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.metalmarket.eu/en/products/krugerrand-1-oz-gold-various-vintages-1786"},
	{Site: "metalmarket", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.metalmarket.eu/en/products/canadian-maple-leaf-1-ounce-of-gold-2026-13635"},
	{Site: "metalmarket", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.metalmarket.eu/en/products/canadian-maple-leaf-1-oz-silver-2026-13642"},
	{Site: "mennicakapitalowa", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://mennicakapitalowa.pl/product-pol-265-Sztabka-zlota-LBMA-1-kg.html"},
	{Site: "mennicakapitalowa", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://mennicakapitalowa.pl/product-pol-443-Moneta-zlota-Krugerrand-2026-1-oz-48h.html"},
	{Site: "mennicakapitalowa", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://mennicakapitalowa.pl/product-pol-227-Moneta-zlota-Kanadyjski-Lisc-Klonowy-2026-1-oz-48h.html"},
	// kupnozlota's JSON-LD holds the EUR price labelled PLN; meta is right.
	{Site: "kupnozlota", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, Sources: metaOnly, URL: "https://www.kupnozlota.pl/1-kg-goldbarren-valcambi.html"},
	{Site: "kupnozlota", Description: "Silver Maple Leaf 1oz (circulated)", WeightGrams: oz, Sources: metaOnly, URL: "https://www.kupnozlota.pl/1-oz-silbermuenze-maple-leaf-zirkuliert.html"},

	// Czech Republic
	{Site: "goldenhouse", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://goldenhouse.cz/de/c/zlato-de/argor-heraeus-investment-goldbarren-1000g"},
	{Site: "goldenhouse", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://goldenhouse.cz/de/c/zlato-de/1-unze-krugerrand-goldmunze"},
	{Site: "goldenhouse", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://goldenhouse.cz/de/c/zlato-de/maple-leaf-1-oz-2026-goldmunze"},
	{Site: "goldenhouse", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://goldenhouse.cz/de/c/stribro-de/maple-leaf-1-oz-silbermunze-2026"},
	{Site: "goldengate", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://eshop.goldengate.cz/products/zlaty-slitek-1000-g"},
	{Site: "goldengate", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://eshop.goldengate.cz/products/krugerrand-1-oz-zlata-mince"},
	{Site: "goldengate", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://eshop.goldengate.cz/products/maple-leaf-1-oz-zlata-mince"},
	{Site: "goldengate", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://eshop.goldengate.cz/products/maple-leaf-1-oz-stribrna-mince"},
	{Site: "auportal", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.auportal.cz/investicni-zlato-zlaty-slitek-1000g-argor-heraeus-sa_z670/"},
	{Site: "auportal", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.auportal.cz/zlata-mince-50-cad-maple-leaf-1-oz_z608/"},
	{Site: "auportal", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.auportal.cz/stribrna-mince-5-cad-maple-leaf-1-oz-2026_z3644/"},

	// Slovakia
	{Site: "zlatypristav", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://zlatypristav.sk/obchod/zlato-gbv/investicne-zlate-mince/maple-leaf-1-oz-gold-2026/"},
	{Site: "eshop-zlato.sk", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.eshop-zlato.sk/zlata-tehlicka-argor-heraeus-1000g"},
	{Site: "slovenske-mince", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://slovenske-mince.sk/argor-heraeus-sa-1000-gramov-investicna-zlata-tehlicka-investice-0000387"},
	{Site: "slovenske-mince", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://slovenske-mince.sk/kruger-rand-1-oz-unc-investicni-zlata-minca-investice-0000338"},
	{Site: "slovenske-mince", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://slovenske-mince.sk/maple-leaf-1-oz-investicni-zlata-minca-investice-0000340"},
	{Site: "slovenske-mince", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://slovenske-mince.sk/maple-leaf-1-oz-unc-investicni-stribrna-mince-investice-0000256"},

	// Hungary (HUF)
	{Site: "tavex.hu", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.hu/arany/1oz-del-afrikai-krugerrand/"},
	{Site: "tavex.hu", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.hu/arany/1oz-kanadai-juharlevel-erme/"},
	{Site: "tavex.hu", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.hu/ezuest/1oz-kanadai-juharlevel-ezuest-erme/"},
	{Site: "aranypiac", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://www.aranypiac.hu/argor-heraeus-munze-osterreich-aranyrud-1000-gramm"},
	{Site: "aranypiac", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.aranypiac.hu/Krugerrand-1-uncia-Del-Afrika"},
	{Site: "aranypiac", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.aranypiac.hu/Juharlevel-1-uncia-Kanada"},
	{Site: "aranypiac", Description: "Silver Maple Leaf 1oz (min. 10)", WeightGrams: oz, URL: "https://www.aranypiac.hu/Juharlevel-EZUST-1/1-uncia-min-rendelheto-25db"},
	// aranykereskedes's JSON-LD holds the EUR price labelled HUF; meta is right.
	{Site: "aranykereskedes", Description: "Gold bar 1kg Heraeus", WeightGrams: kg, Sources: metaOnly, URL: "https://aranykereskedes.hu/1-kg-heraeus-aranyrud.html"},
	{Site: "aranykereskedes", Description: "Gold Krugerrand 1oz", WeightGrams: oz, Sources: metaOnly, URL: "https://aranykereskedes.hu/1-uncia-krugerrand-aranyerme.html"},
	{Site: "aranykereskedes", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, Sources: metaOnly, URL: "https://aranykereskedes.hu/1uncia-maple-leaf-aranyerme.html"},
	{Site: "aranykereskedes", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, Sources: metaOnly, URL: "https://aranykereskedes.hu/1-uncia-maple-leaf-ezusterme-2022.html"},

	// Romania (RON)
	{Site: "tavex.ro", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://tavex.ro/aur/1000gr-lingou-dinaur/"},
	{Site: "tavex.ro", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.ro/aur/krugerrand-sud-african-1uncie/"},
	{Site: "tavex.ro", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.ro/aur/moneda-din-aur-frunza-artar-canada-1uncie/"},
	{Site: "tavex.ro", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.ro/argint/moneda-din-argint-1uncie-frunza-de-artar/"},
	{Site: "goldbars.ro", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://goldbars.ro/catalog/lingou-aur-1kg-argor-heraeus-24k"},
	{Site: "goldbars.ro", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://goldbars.ro/catalog/moneda-aur-1oz-31-10g-canada-maple-24k"},
	{Site: "goldbars.ro", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://goldbars.ro/catalog/moneda-argint-1oz-31-10g-canada-maple-2026-999"},
	{Site: "avangardgold", Description: "Gold bar 1kg Metalor", WeightGrams: kg, URL: "https://avangardgold.ro/products/lingou-1-kg-metalor"},
	{Site: "avangardgold", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://avangardgold.ro/products/moneda-aur-krugerrand"},
	{Site: "avangardgold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://avangardgold.ro/products/moneda-de-aur-canada-1-uncie"},
	{Site: "avangardgold", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://avangardgold.ro/products/moneda-de-argint-1-uncie-frunza-de-ar%C8%9Bar-canada"},
	{Site: "neogold", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, URL: "https://www.neogold.ro/produs/lingou-de-aur-valcambi-suisse-1kg/"},
	{Site: "neogold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.neogold.ro/produs/moneda-din-aur-royal-canadian-mint-maple-leaf-1-oz-31-1-gr/"},
	{Site: "neogold", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.neogold.ro/produs/moneda-din-argint-maple-leaf-2023-1-oz-31-10-gr-disponibil-in-stoc/"},
	{Site: "aurom", Description: "Gold bar 1kg Argor-Heraeus", WeightGrams: kg, URL: "https://www.aurominvestment.ro/produs/1000-grame-lingou-de-aur-argor-heraeus/"},
	{Site: "aurom", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.aurominvestment.ro/produs/moneda-aur-1oz-krugerrand/"},
	{Site: "aurom", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.aurominvestment.ro/produs/1-uncie-moneda-de-argint-canada-maple-leaf/"},
	{Site: "auracasa", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://auracasa.ro/ro/1-oz/1847-maple-leaf-2026-1-oz-moneda-de-argint-pentru-investitii.html?SubmitCurrency=1&id_currency=2"},

	// Bulgaria
	{Site: "tavex.bg", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://tavex.bg/en/gold/1000g-gold-bar/"},
	{Site: "tavex.bg", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.bg/en/gold/1oz-south-african-krugerrand/"},
	{Site: "tavex.bg", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.bg/en/gold/1oz-canadian-maple-leaf-gold-coin/"},
	{Site: "tavex.bg", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.bg/en/silver/1oz-canadian-maple-leaf-silver-coin/"},
	{Site: "zlatenrezerv", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, URL: "https://www.zlatenrezerv.bg/investicionno-zlato/zlatni-kyulcheta/golemi-zlatni-kyulcheta/valcambi-1-kilogram-zlatno-kyulche/"},
	{Site: "zlatenrezerv", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.zlatenrezerv.bg/investicionno-zlato/zlatni-moneti/yuzhnoafrikanski-zlatni-moneti-krugerrand/1-unciya-zlatna-moneta-kryugerrand-yuzhna-afrika/"},
	{Site: "zlatenrezerv", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.zlatenrezerv.bg/investicionno-srebro/srebarni-moneti/kanadski-klenov-list/1-unciya-srebarna-moneta-kanadski-klenov-list-s-dds-na-marzha/"},
	{Site: "topgold", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://topgold.bg/product/1-uncziya-31-1-gr-zlatna-moneta-krugerrand-mixed-years/"},
	{Site: "topgold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://topgold.bg/product/1-uncziya-31-1-gr-zlatna-moneta-kanadski-klenov-list-razlichni-godini/"},
	{Site: "mygold.bg", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://mygold.bg/1-unciya-krugerrand-razlichni-godini-yuzhnoafrikanska-zlatna-moneta"},
	{Site: "mygold.bg", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://mygold.bg/1-unciya-klenov-list-razlichni-godini-kanadska-zlatna-moneta"},
	{Site: "mygold.bg", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://mygold.bg/1-unciya-klenov-list-razlichni-godini-kanadska-srebrna-moneta"},
	{Site: "factorin", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://factorin.bg/bg/product/store/%D0%BC%D0%BE%D0%B4%D0%B5%D1%80%D0%BD%D0%B8-%D0%BC%D0%BE%D0%BD%D0%B5%D1%82%D0%B8/kru31/"},
	{Site: "factorin", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://factorin.bg/bg/product/store/%D0%BC%D0%BE%D0%B4%D0%B5%D1%80%D0%BD%D0%B8-%D0%BC%D0%BE%D0%BD%D0%B5%D1%82%D0%B8/map31/"},
	{Site: "currencycenter", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://currencycenter.eu/investment-gold-coins/zlatna-moneta-krugerrand-yuzhna-afrika?locale=en"},
	{Site: "currencycenter", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://currencycenter.eu/investment-gold-coins/zlatna-moneta-kanadski-klenov-list-1-uncia?locale=en"},

	// Estonia
	{Site: "tavid", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, URL: "https://tavid.ee/kuld/1-kg-valcambi-suisse-valatud-kuldplaat/"},
	{Site: "tavid", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavid.ee/kuld/1oz-louna-aafrika-krugerrand-kuldmunt/"},
	{Site: "tavid", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavid.ee/kuld/1oz-kanada-maple-leaf-kuldmunt/"},
	{Site: "tavid", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavid.ee/hobe/1-oz-kanada-maple-leaf-hobemunt/"},
	{Site: "gvs.ee", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, URL: "https://gvs.ee/1-kg-goldbarren-valcambi.html"},
	{Site: "gvs.ee", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://gvs.ee/1-oz-goldcoin-krugerrand-various.html"},
	{Site: "gvs.ee", Description: "Gold Maple Leaf 1oz (circulated)", WeightGrams: oz, URL: "https://gvs.ee/1-oz-goldcoin-maple-leaf-circulated.html"},
	{Site: "gvs.ee", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://gvs.ee/1-oz-silver-maple-leaf.html"},
	{Site: "eurex.ee", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://eurex.ee/en/shop/1-oz-krugerrand-gold-coin-11607"},
	{Site: "eurex.ee", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://eurex.ee/en/shop/1-oz-maple-leaf-gold-coin-2026-104320"},
	{Site: "eurex.ee", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://eurex.ee/en/shop/1-oz-maple-leaf-silver-coin-2026-204276"},
	{Site: "tavast", Description: "Gold bar 1kg UBS", WeightGrams: kg, URL: "https://gold.tavast.eu/product/kuldplaat-ubs-1000g/"},
	{Site: "tavast", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://gold.tavast.eu/product/hobemunt-maple-leaf-1-oz-2016/"},

	// Latvia
	{Site: "tavex.lv", Description: "Gold bar 1kg", WeightGrams: kg, URL: "https://tavex.lv/en/gold/1000g-gold-bar/"},
	{Site: "tavex.lv", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.lv/en/gold/1oz-south-african-krugerrand/"},
	{Site: "tavex.lv", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.lv/en/gold/1oz-canadian-maple-leaf-gold-coin/"},
	{Site: "tavex.lv", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.lv/en/silver/1oz-canadian-maple-leaf-silver-coin/"},

	// Lithuania
	{Site: "tavex.lt", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://tavex.lt/en/gold/1oz-south-african-krugerrand/"},
	{Site: "tavex.lt", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.lt/en/gold/1oz-canadian-maple-leaf-gold-coin/"},
	{Site: "tavex.lt", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://tavex.lt/en/silver/1oz-canadian-maple-leaf-silver-coin/"},
	{Site: "eurodata.lt", Description: "Gold bar 1kg Valcambi", WeightGrams: kg, URL: "https://www.eurodata.lt/1-kg-investicinio-aukso-luitas-valcambi-999-9-lietas"},
	{Site: "eurodata.lt", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.eurodata.lt/1-oz-31-10-g-auksine-moneta-krugerrand-pietu-afrikos-respublika-2026"},
	{Site: "eurodata.lt", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.eurodata.lt/1-oz-31-10-g-auksine-moneta-klevo-lapas-kanada-2026"},
	{Site: "eurodata.lt", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.eurodata.lt/1-oz-31-10-g-sidabrine-moneta-klevo-lapas-kanada-2026"},
	{Site: "gerettigold", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://gerettigold.lt/parduotuve/auksines-monetos/1-oz-klevo-lapas-auksine-moneta-2024/"},
	{Site: "gerettigold", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://gerettigold.lt/parduotuve/sidabrines-monetos/1-oz-klevo-lapas-sidabrine-moneta/"},
	{Site: "monetupasaulis", Description: "Gold Krugerrand 1oz", WeightGrams: oz, URL: "https://www.monetupasaulis.lt/preke/1-oz-auksine-moneta-krugerrand-2020-pietu-afrika/"},
	{Site: "monetupasaulis", Description: "Gold Maple Leaf 1oz", WeightGrams: oz, URL: "https://www.monetupasaulis.lt/preke/1-oz-auksine-moneta-kanados-klevo-lapas-2026-kanada/"},
	{Site: "aurea.shop", Description: "Silver Maple Leaf 1oz", WeightGrams: oz, URL: "https://aurea.shop/taupymui-tinkancios-bullions/sidabrine-moneta-kanados-klevo-lapas-1-oz-2026.html"},
}

var (
	proaurum   = []scraper.Source{scraper.Meta}
	metaOnly   = []scraper.Source{scraper.Meta}
	cssOnly    = []scraper.Source{scraper.CSS}
	proaurumCH = []string{".price-ask_price .price"}
	// goldavenue defaults to CHF.
	eur      = map[string]string{"currency": "EUR"}
	sgbPrice = []string{
		`div[class~="tw:grid-cols-10"] > div:containsOwn("1+") + div`,
		`h1 + div > div[class~="tw:font-bold"]`,
	}
)
