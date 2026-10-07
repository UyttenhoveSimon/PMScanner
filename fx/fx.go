// Package fx converts prices to euros using the European Central Bank's
// daily reference rates, completed by open.er-api.com for currencies the ECB
// does not publish (e.g. RUB, KZT, AZN, RSD).
package fx

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const (
	ecbURL      = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	fallbackURL = "https://open.er-api.com/v6/latest/EUR"
)

var client = &http.Client{Timeout: 15 * time.Second}

// FetchRates returns the ECB rates, plus fallback rates for the currencies
// the ECB does not cover. It fails only when neither source answers.
func FetchRates() (Rates, error) {
	rates, ecbErr := FetchECB()
	extra, fallbackErr := FetchFallback()
	if ecbErr != nil && fallbackErr != nil {
		return nil, errors.Join(ecbErr, fallbackErr)
	}
	if rates == nil {
		rates = Rates{}
	}
	for cur, r := range extra {
		if _, ok := rates[cur]; !ok && cur != "EUR" {
			rates[cur] = r
		}
	}
	return rates, nil
}

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

// FetchFallback downloads rates from open.er-api.com, a free daily source
// covering most world currencies.
func FetchFallback() (Rates, error) {
	resp, err := client.Get(fallbackURL)
	if err != nil {
		return nil, fmt.Errorf("fetch fallback rates: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch fallback rates: HTTP %d", resp.StatusCode)
	}
	return parseFallback(resp.Body)
}

func parseFallback(r io.Reader) (Rates, error) {
	var body struct {
		Result string                 `json:"result"`
		Rates  map[string]json.Number `json:"rates"`
	}
	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return nil, fmt.Errorf("parse fallback rates: %w", err)
	}
	if body.Result != "success" || len(body.Rates) == 0 {
		return nil, fmt.Errorf("parse fallback rates: no rates (result %q)", body.Result)
	}
	rates := Rates{}
	for cur, n := range body.Rates {
		if v, err := decimal.NewFromString(n.String()); err == nil && v.IsPositive() {
			rates[cur] = v
		}
	}
	return rates, nil
}
