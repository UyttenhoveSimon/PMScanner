package websites

import (
	"pmscanner/types"

	"github.com/gocolly/colly/v2"
)

func collect(c *colly.Collector) []types.Item {
	var url = "https://or.fr/"
	var item = types.Item{}
	var item2 = types.Item{}
	var items = []types.Item{item, item2}

	return items
}
