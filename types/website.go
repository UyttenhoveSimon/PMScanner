package types

import "github.com/gocolly/colly/v2"

type Website interface {
	Collect(c *colly.Collector) []Item
}
