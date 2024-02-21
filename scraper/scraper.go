package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/khdiyz/web-scraper/utils"
)

type Device struct {
	Title         string `json:"title"`
	StarAmount    string `json:"star"`
	Comment       string `json:"numberOfComment"`
	OldPrice      string `json:"oldPrice"`
	Price         string `json:"price"`
	PricePerMonth string `json:"pricePerMonth"`
	ImageUrl      string `json:"image"`
}

func ScrapeDevices() (devices []Device, err error) {
	c := colly.NewCollector(
		colly.AllowedDomains("elmakon.uz"),
	)

	c.OnHTML("div.ut2-gl__body", func(h *colly.HTMLElement) {
		device := Device{
			Title:         h.ChildText("div.ut2-gl__name"),
			StarAmount:    utils.TrimSpacesLR(h.ChildText("a.ty-product-review-reviews-stars__link")),
			Comment:       utils.TrimSpacesLR(h.ChildText("div.cn-reviews")),
			OldPrice:      utils.TrimSpacesLR(h.ChildText("span.ty-strike")),
			Price:         utils.TrimSpacesLR(h.ChildText("span.ty-price-num")),
			PricePerMonth: utils.TrimSpacesLR(h.ChildText("p")),
			ImageUrl:      h.ChildAttr("img", "src"),
		}
		devices = append(devices, device)
	})

	c.OnHTML("a.ty-pagination__next", func(h *colly.HTMLElement) {
		nextPage := h.Attr("href")
		if nextPage != "" {
			c.Visit(h.Request.AbsoluteURL(nextPage))
		}
	})

	err = c.Visit("https://elmakon.uz/telefony-gadzhety-aksessuary/")
	if err != nil {
		return nil, err
	}

	c.Wait()

	return devices, nil
}
