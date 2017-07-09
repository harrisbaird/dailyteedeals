package v2

import (
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/modext"
)

type v2Site struct {
	ID       int          `json:"id"`
	Slug     string       `json:"slug"`
	Name     string       `json:"name"`
	Products []*v2Product `json:"products"`
}

type v2Artist struct {
	ID      int         `json:"id"`
	Slug    string      `json:"slug"`
	Name    string      `json:"name"`
	Designs []*v2Design `json:"designs"`
}

type v2Design struct {
	ID       int            `json:"id"`
	Slug     string         `json:"slug"`
	Name     string         `json:"name"`
	URL      string         `json:"url"`
	Artist   *v2Artist      `json:"artist"`
	Products []*v2Product   `json:"products"`
	Order    map[string]int `json:"order"`
}

type v2Price struct {
	ID string `json:"id"`
}

type v2Images struct {
	Image300        string `json:"small"`
	Image1200       string `json:"large"`
	BackgroundColor string `json:"backgroundColor"`
}

type v2Product struct {
	ID         int                         `json:"id"`
	Slug       string                      `json:"slug"`
	BuyURL     string                      `json:"buyURL"`
	Active     bool                        `json:"active"`
	Deal       bool                        `json:"deal"`
	LastChance bool                        `json:"lastChance"`
	Site       *v2Site                     `json:"site"`
	Prices     map[string]ApproximatePrice `json:"prices"`
	Images     *v2Images                   `json:"images"`
}

func buildV2Product(product *models.Product) *v2Product {
	var prices map[string]ApproximatePrice
	if product.Prices != nil {
		prices = ConvertPrices(product.Prices)
	}

	return &v2Product{
		ID:         product.ID,
		Slug:       product.Slug,
		BuyURL:     modext.ProductGoURL(product),
		Active:     product.Active,
		Deal:       product.Deal,
		LastChance: product.LastChance,
		Site:       buildV2Site(product.R.Site),
		Prices:     prices,
		Images:     buildV2Images(product),
	}
}

func buildV2Products(products models.ProductSlice) []*v2Product {
	var out []*v2Product

	if products == nil {
		return []*v2Product{}
	}

	for _, product := range products {
		out = append(out, buildV2Product(product))
	}

	return out
}

func buildV2Images(product *models.Product) *v2Images {
	return &v2Images{
		Image300:        modext.ProductImageURL(product, 300),
		Image1200:       modext.ProductImageURL(product, 1200),
		BackgroundColor: product.ImageBackground,
	}
}

func buildV2Artist(artist *models.Artist) *v2Artist {
	out := v2Artist{
		ID:      artist.ID,
		Slug:    artist.Slug,
		Name:    artist.Name,
		Designs: buildV2Designs(artist.R.Designs),
	}

	return &out
}

func buildV2ArtistWithDesigns(artist *models.Artist) *v2Artist {
	out := buildV2Artist(artist)
	for _, d := range artist.R.Designs {
		out.Designs = append(out.Designs, buildV2Design(d))
	}
	return out
}

func buildV2Design(design *models.Design) *v2Design {
	out := v2Design{
		ID:       design.ID,
		Slug:     design.Slug,
		Name:     design.Name,
		URL:      "https://dailyteedeals.com/designs/" + design.Slug,
		Artist:   buildV2Artist(design.R.Artist),
		Products: buildV2Products(design.R.Products),
	}
	return &out
}

func buildV2Designs(designs models.DesignSlice) []*v2Design {
	if designs == nil {
		return []*v2Design{}
	}
	var out []*v2Design

	for _, design := range designs {
		out = append(out, buildV2Design(design))
	}

	return out
}

func buildV2Site(site *models.Site) *v2Site {
	out := v2Site{
		ID:       site.ID,
		Slug:     site.Slug,
		Name:     site.Name,
		Products: []*v2Product{},
	}

	return &out
}

func buildV2SiteWithProducts(site *models.Site) *v2Site {
	out := buildV2Site(site)
	out.Products = buildV2Products(site.R.Products)
	return out
}

func buildV2Sites(sites models.SiteSlice) []*v2Site {
	var out []*v2Site

	for _, site := range sites {
		out = append(out, buildV2Site(site))
	}

	return out
}
