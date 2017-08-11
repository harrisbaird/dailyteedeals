package server

import (
	"fmt"

	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/utils"
)

type v2Site struct {
	ID       int          `json:"id"`
	Slug     string       `json:"slug"`
	Name     string       `json:"name"`
	Products []*v2Product `json:"products,omitempty"`
}

type v2Artist struct {
	ID      int         `json:"id"`
	Slug    string      `json:"slug"`
	Name    string      `json:"name"`
	Designs []*v2Design `json:"designs,omitempty"`
}

type v2Design struct {
	ID         int           `json:"id"`
	Slug       string        `json:"slug"`
	Name       string        `json:"name"`
	URL        string        `json:"url"`
	Artist     *v2Artist     `json:"artist"`
	Products   []*v2Product  `json:"products,omitempty"`
	Categories []*v2Category `json:"categories"`
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
	ID         int                               `json:"id"`
	Design     *v2Design                         `json:"design,omitempty"`
	Site       *v2Site                           `json:"site"`
	Slug       string                            `json:"slug"`
	BuyURL     string                            `json:"buyURL"`
	Active     bool                              `json:"active"`
	Deal       bool                              `json:"deal"`
	LastChance bool                              `json:"lastChance"`
	Prices     map[string]utils.ApproximatePrice `json:"prices"`
	Images     *v2Images                         `json:"images"`
}

type v2Category struct {
	ID       int    `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func buildV2Product(product *models.Product) *v2Product {
	return &v2Product{
		ID:         product.ID,
		Design:     buildV2Design(product.Design),
		Site:       buildV2Site(product.Site),
		Slug:       product.Slug,
		BuyURL:     product.GoURL(),
		Active:     product.Active,
		Deal:       product.Deal,
		LastChance: product.LastChance,
		Prices:     product.ConvertedPrices,
		Images:     buildV2Images(product),
	}
}

func buildV2Products(products []*models.Product) []*v2Product {
	out := []*v2Product{}

	for _, product := range products {
		out = append(out, buildV2Product(product))
	}

	return out
}

func buildV2Images(product *models.Product) *v2Images {
	return &v2Images{
		Image300:        product.ImageURL(models.SmallImageType),
		Image1200:       product.ImageURL(models.LargeImageType),
		BackgroundColor: product.ImageBackground,
	}
}

func buildV2Artist(artist *models.Artist) *v2Artist {
	return &v2Artist{
		ID:      artist.ID,
		Slug:    artist.Slug,
		Name:    artist.Name,
		Designs: []*v2Design{},
	}
}

func buildV2ArtistWithDesigns(artist *models.Artist) *v2Artist {
	out := buildV2Artist(artist)
	out.Designs = buildV2DesignsWithProducts(artist.Designs)
	return out
}

func buildV2Design(design *models.Design) *v2Design {
	if len(design.Categories) > 0 {
		fmt.Printf("%d\n", design.ID)
	}

	return &v2Design{
		ID:         design.ID,
		Slug:       design.Slug,
		Name:       design.Name,
		URL:        "https://dailyteedeals.com/designs/" + design.Slug,
		Artist:     buildV2Artist(design.Artist),
		Categories: buildV2Categories(design.Categories),
		Products:   []*v2Product{},
	}
}

func buildV2DesignWithProducts(design *models.Design) *v2Design {
	out := buildV2Design(design)
	out.Products = buildV2Products(design.Products)
	return out
}

func buildV2DesignsWithProducts(designs []*models.Design) []*v2Design {
	out := []*v2Design{}
	for _, design := range designs {
		out = append(out, buildV2DesignWithProducts(design))
	}
	return out
}

func buildV2Designs(designs []*models.Design) []*v2Design {
	out := []*v2Design{}
	if designs == nil || len(designs) == 0 {
		return out
	}

	for _, design := range designs {
		out = append(out, buildV2Design(design))
	}

	return out
}

func buildV2Site(site *models.Site) *v2Site {
	return &v2Site{
		ID:       site.ID,
		Slug:     site.Slug,
		Name:     site.Name,
		Products: []*v2Product{},
	}
}

func buildV2SiteWithProducts(site *models.Site) *v2Site {
	out := buildV2Site(site)
	out.Products = buildV2Products(site.Products)
	return out
}

func buildV2Sites(sites []*models.Site) []*v2Site {
	out := []*v2Site{}

	for _, site := range sites {
		out = append(out, buildV2Site(site))
	}

	return out
}

func buildV2Category(category *models.Category) *v2Category {
	return &v2Category{
		ID:       category.ID,
		Slug:     category.Slug,
		Name:     category.Name,
		ImageURL: category.Product.LargeImageURL(),
	}
}

func buildV2Categories(categories []*models.Category) []*v2Category {
	out := []*v2Category{}

	for _, category := range categories {
		out = append(out, buildV2Category(category))
	}

	return out
}
