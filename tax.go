package main

import (
	"github.com/shopspring/decimal"

	"pmscanner/scraper"
)

// countryTax summarises how precious metals are taxed in a country, for the
// "Taxes by country" section and to add VAT to prices quoted without it.
// Rates are standard VAT rates as of 2026; they are general information,
// not tax advice.
type countryTax struct {
	VAT    decimal.Decimal // standard rate, in percent
	Gold   string          // investment gold
	Others string          // silver, platinum, palladium, copper…
}

const (
	euGold     = "Exempt"
	euOthers   = "Standard VAT, or margin scheme on some coins."
	notSpecial = "Standard VAT."
)

func rate(s string) decimal.Decimal { return decimal.RequireFromString(s) }

var taxes = map[string]countryTax{
	"AT": {rate("20"), euGold, euOthers},
	"BE": {rate("21"), euGold, euOthers},
	"BG": {rate("20"), euGold, euOthers},
	"CY": {rate("19"), euGold, euOthers},
	"CZ": {rate("21"), euGold, euOthers},
	"DE": {rate("19"), euGold, euOthers + " Margin-taxed coins are labelled \"differenzbesteuert\"."},
	"DK": {rate("25"), euGold, euOthers},
	"EE": {rate("24"), euGold, euOthers},
	"ES": {rate("21"), euGold, euOthers},
	"FI": {rate("25.5"), euGold, euOthers},
	"FR": {rate("20"), euGold, euOthers + " Metal kept in a dealer's vault (e.g. or.fr) is VAT-free until delivered."},
	"GR": {rate("24"), euGold, euOthers},
	"HR": {rate("25"), euGold, euOthers},
	"HU": {rate("27"), euGold, euOthers},
	"IE": {rate("23"), euGold, euOthers},
	"IT": {rate("22"), euGold, euOthers},
	"LT": {rate("21"), euGold, euOthers},
	"LV": {rate("21"), euGold, euOthers},
	"NL": {rate("21"), euGold, euOthers},
	"PL": {rate("23"), euGold, euOthers},
	"PT": {rate("23"), euGold, euOthers},
	"RO": {rate("21"), euGold, euOthers},
	"SE": {rate("25"), euGold, euOthers},
	"SI": {rate("22"), euGold, euOthers},
	"SK": {rate("23"), euGold, euOthers},
	"CH": {rate("8.1"), "Exempt", "Swiss VAT. Metal kept in a bonded warehouse or free port is VAT-free until delivered."},
	"LI": {rate("8.1"), "Exempt", "Swiss VAT (Liechtenstein is in the Swiss VAT area)."},
	"GB": {rate("20"), "Exempt", "20% VAT. Second-hand coins may be sold under the margin scheme."},
	"NO": {rate("25"), "Exempt", notSpecial},
	"TR": {rate("20"), "Bullion generally exempt", "Bullion deliveries are generally exempt; other items carry VAT."},
	"RU": {rate("22"), "Exempt for bars and investment coins bought by individuals", "Bars bought by individuals are exempt; other items carry VAT."},
	"KZ": {rate("16"), "Bars from the National Bank are exempt", notSpecial},
	"AZ": {rate("18"), "Not documented here", notSpecial},
}

// withVAT returns a price per gram including VAT, adding the country's
// standard rate when the shop quotes it without VAT.
func withVAT(perGram decimal.Decimal, vat scraper.VAT, country string) (decimal.Decimal, bool) {
	if vat != scraper.VATExcluded {
		return perGram, true
	}
	t, ok := taxes[country]
	if !ok {
		return perGram, false
	}
	return perGram.Mul(decimal.NewFromInt(1).Add(t.VAT.Div(decimal.NewFromInt(100)))), true
}
