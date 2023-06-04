package types

import "github.com/shopspring/decimal"

type Item struct {
	Price       decimal.Decimal
	Currency    string
	Description string
	Weight      string
	Page        string
	HtmlFilter  string
}
