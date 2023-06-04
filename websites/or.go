package websites

import (
	"fmt"
	. "pmscanner/types"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
)

type Or struct{}

var itemsOr = []Item{
	{
		Page:        "https://or.fr/produits/acheter-or/lingots-or/lingot-d-or-1-kilogramme-valcambi-60",
		HtmlFilter:  "li.text-light.mb-4",
		Weight:      "1kg",
		Description: "ValCambi 1kg",
		Currency:    "€",
	},
}

func (o Or) Collect(c *colly.Collector) []Item {

	var cp = c.Clone()
	for _, element := range itemsOr {

		cp.OnHTML(element.HtmlFilter, func(e *colly.HTMLElement) {
			priceFormat := regexp.MustCompile(`(\d+.\d+\,\d+)`)
			price := priceFormat.FindString(e.Text)
			fmt.Printf("found: %s\n", price)
			twPrice := strings.ReplaceAll(price, ",", ".")
			twPrice = strings.ReplaceAll(twPrice, " ", "")
			fmt.Printf("tweaked Price: %s\n", twPrice)
			element.Price, _ = decimal.NewFromString(twPrice)
			fmt.Printf("stored price: %s\n", element.Price)
		})
		cp.Visit(element.Page)
	}

	return itemsOr
}
