package modext

import (
	"database/sql"
	"net/url"
	"strings"

	"log"

	titlecase "github.com/AlasdairF/Titlecase"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

func init() {
	models.AddArtistHook(boil.BeforeInsertHook, artistSaveHook)
	models.AddArtistHook(boil.BeforeUpdateHook, artistSaveHook)
	models.AddArtistHook(boil.BeforeUpsertHook, artistSaveHook)
}

// FindOrCreateArtist convenience function.
// It attempts to find an artist by name or by urls.
// Otherwise it creates a new artist.
func FindOrCreateArtist(db boil.Executor, name string, urls []string) (*models.Artist, error) {
	artist, err := models.Artists(db,
		qm.Where(TagQuery("urls", urls, NormalizeUrls)),
		qm.Or(TagQuery("tags", []string{name}, NormalizeTags)),
	).One()

	if err == sql.ErrNoRows {
		artist = &models.Artist{Name: name}
		ArtistAppendWhitelistedUrls(artist, urls)
		err = artist.Insert(db)
	}

	return artist, err
}

func FindArtistBySlug(db boil.Executor, slug string) (*models.Artist, error) {
	return models.Artists(db, qm.Where("slug=?", slug)).One()
}

var whitelistDomains = []string{"twitter.com", "facebook.com", "redbubble.com",
	"neatoshop.com", "teepublic.com", "instagram.com", "society6.com"}

// ArtistAppendWhitelistedUrls normalizes urls and merges any urls from
// whitelisted domains with the artists urls.
func ArtistAppendWhitelistedUrls(a *models.Artist, urls []string) {
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

func artistSaveHook(exec boil.Executor, a *models.Artist) error {
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
