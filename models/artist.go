package models

import (
	"log"
	"net/url"
	"strings"

	titlecase "github.com/AlasdairF/Titlecase"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
)

var whitelistDomains = []string{"twitter.com", "facebook.com", "redbubble.com",
	"neatoshop.com", "teepublic.com", "instagram.com", "society6.com"}

type Artist struct {
	ID      int
	Name    string
	Slug    string
	Urls    []string `pg:",array" sql:",notnull"`
	Tags    []string `pg:",array" sql:",notnull"`
	Designs []*Design
}

// FindOrCreateArtist convenience function.
// It attempts to find an artist by name or by urls.
// Otherwise it creates a new artist.
func FindOrCreateArtist(db orm.DB, name string, urls []string) (*Artist, error) {
	var artist Artist
	err := db.Model(&artist).
		Where(TagQuery("urls", urls, NormalizeUrls)).
		WhereOr(TagQuery("tags", []string{name}, NormalizeTags)).
		First()

	if err == pg.ErrNoRows {
		artist = Artist{Name: name}
		artist.ArtistAppendWhitelistedUrls(urls)
		err = db.Insert(&artist)
	}

	return &artist, err
}

func FindArtistBySlug(db orm.DB, slug string, page int) (*Artist, error) {
	perPage := config.App.ItemsPerPage

	var artist Artist
	err := db.Model(&artist).
		Column("artist.*", "Designs", "Designs.Artist", "Designs.Products", "Designs.Products.Site").
		Relation("Designs", func(q *orm.Query) (*orm.Query, error) {
			return q.Offset(perPage * (page - 1)).Limit(perPage), nil
		}).
		Where("artist.slug=?", slug).
		First()
	return &artist, err
}

// func (a *Artist) LoadActiveDesigns(db orm.DB, page int) ([]*Design, error) {
// 	var designs []*Design
// 	err := db.Model(&designs).
// 		Column("design.*", "Artist", "Products", "Products.Site").
// 		Join("LEFT OUTER JOIN products AS pj ON design.id = pj.design_id").
// 		// Group("design.id", "artist.id", "pj.active").
// 		// Where("pj.active=?", true).
// 		Where("artist_id=?", a.ID).
// 		Relation("Products", func(q *orm.Query) (*orm.Query, error) {
// 			return q.Where("product.active=?", true).
// 				Order("product.deal DESC").
// 				OrderExpr("product.prices -> 'usd' ASC"), nil
// 		}).
// 		Order("design.name").
// 		Offset(config.App.ItemsPerPage * (page - 1)).
// 		Limit(config.App.ItemsPerPage).
// 		Select()
// 	return designs, err
// }

// ArtistAppendWhitelistedUrls normalizes urls and merges any urls from
// whitelisted domains with the artists urls.
func (a *Artist) ArtistAppendWhitelistedUrls(urls []string) {
	urls = NormalizeUrls(urls)
	var output []string
	for _, u := range urls {
		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Printf("Unable to parse URL: %s - %v", u, err)
			continue
		}

		for _, domain := range whitelistDomains {
			if parsedURL.Hostname() == domain {
				output = append(output, u)
			}
		}
	}

	a.Urls = append(a.Urls, output...)
}

func (a *Artist) BeforeInsert(db orm.DB) error {
	return a.normalize(db)
}

func (a *Artist) BeforeUpdate(db orm.DB) error {
	return a.normalize(db)
}

func (a *Artist) normalize(db orm.DB) error {
	name := titlecase.English(a.Name)
	a.Name = strings.TrimSpace(name)
	if a.Name != "" {
		a.Tags = append(a.Tags, a.Name)
	}
	if a.Slug == "" {
		a.Slug = MakeSlug(a.Name)
	}
	a.Tags = NormalizeTags(a.Tags)
	a.Urls = NormalizeUrls(a.Urls)
	return nil
}
