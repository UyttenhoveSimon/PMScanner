// Package fx converts prices to euros using the European Central Bank's
// daily reference rates.
package fx

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

// Rates maps a currency code to its value of one euro, e.g. "GBP": 0.85.
type Rates map[string]decimal.Decimal

// ToEUR converts an amount to euros. It fails for unknown currencies; an empty
// currency is assumed to be EUR, the default for the shops scanned.
func (r Rates) ToEUR(amount decimal.Decimal, currency string) (decimal.Decimal, error) {
	if currency == "" || currency == "EUR" {
		return amount, nil
	}
	rate, ok := r[currency]
	if !ok || !rate.IsPositive() {
		return decimal.Zero, fmt.Errorf("no exchange rate for %q", currency)
	}
	return amount.Div(rate), nil
}

// FetchECB downloads today's ECB reference rates.
func FetchECB() (Rates, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(ecbURL)
	if err != nil {
		return nil, fmt.Errorf("fetch ECB rates: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch ECB rates: HTTP %d", resp.StatusCode)
	}
	return parseECB(resp.Body)
}

type ecbEnvelope struct {
	Rates []struct {
		Currency string `xml:"currency,attr"`
		Rate     string `xml:"rate,attr"`
	} `xml:"Cube>Cube>Cube"`
}

func parseECB(r io.Reader) (Rates, error) {
	var env ecbEnvelope
	if err := xml.NewDecoder(r).Decode(&env); err != nil {
		return nil, fmt.Errorf("parse ECB rates: %w", err)
	}
	rates := Rates{}
	for _, c := range env.Rates {
		if v, err := decimal.NewFromString(c.Rate); err == nil {
			rates[c.Currency] = v
		}
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("parse ECB rates: no rates found")
	}
	return rates, nil
}
