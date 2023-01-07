package websites

import (
	"fmt"
	. "pmscanner/types"

	"github.com/gocolly/colly/v2"
)

type Or struct{}

var itemsOr = []Item{
	{
		Page:        "https://or.fr/produits/acheter-or/lingots-or/lingot-d-or-1-kilogramme-valcambi-60",
		HtmlFilter:  "price text-light me-3",
		Weigth:      "1kg",
		Description: "ValCambi 1kg",
	},
}

func (o Or) collect(c *colly.Collector) []Item {

	var cp = c.Clone()
	for _, element := range itemsOr {

		cp.OnHTML(element.HtmlFilter, func(e *colly.HTMLElement) {
			fmt.Printf("found: %s\n", e.Text)

		})
		cp.Visit(element.Page)
	}

	return itemsOr
}
