package v1

import (
	"fmt"
	"math"
	"strings"

	"github.com/harrisbaird/dailyteedeals/api/v2"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/modext"
)

type v1Site struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type v1Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type v1Design struct {
	ID     int       `json:"id"`
	Name   string    `json:"name"`
	Artist *v1Artist `json:"artist"`
}

type v1Image struct {
	URL             string `json:"url"`
	BackgroundColor string `json:"background_color"`
}

type v1Product struct {
	ID           int       `json:"id"`
	BuyURL       string    `json:"buy_url"`
	DisplayPrice string    `json:"display_price"`
	LastChance   bool      `json:"last_chance"`
	Image        *v1Image  `json:"image"`
	Design       *v1Design `json:"design"`
	Site         *v1Site   `json:"site"`
}

func buildV1Api(products models.ProductSlice) []v1Product {
	var out []v1Product

	for _, product := range products {
		out = append(out, v1Product{
			ID:           product.ID,
			BuyURL:       modext.ProductBuyURL(product),
			DisplayPrice: buildV1Price(product),
			LastChance:   product.LastChance,
			Design:       buildV1Design(product.R.Design),
			Site:         buildV1Site(product.R.Site),
			Image: &v1Image{
				URL:             modext.ProductImageURL(product, 300),
				BackgroundColor: product.ImageBackground,
			},
		})
	}

	return out
}

var v1Currencies = []string{"USD", "GBP", "EUR"}

func buildV1Price(product *models.Product) string {
	prices := v2.ConvertPrices(product.Prices)

	output := []string{}
	for _, currency := range v1Currencies {
		price := prices[currency]
		symbol := v2.CurrencySymbols[price.Currency]
		amount := math.Ceil(float64(price.Amount) / 100)
		output = append(output, fmt.Sprintf("%s%.2f", symbol, amount))
	}

	return strings.Join(output, " / ")
}

func buildV1Artist(artist *models.Artist) *v1Artist {
	return &v1Artist{
		ID:   artist.ID,
		Name: artist.Name,
	}
}

func buildV1Design(design *models.Design) *v1Design {
	return &v1Design{
		ID:     design.ID,
		Name:   design.Name,
		Artist: buildV1Artist(design.R.Artist),
	}
}

func buildV1Site(site *models.Site) *v1Site {
	return &v1Site{
		ID:   site.ID,
		Name: site.Name,
	}
}
