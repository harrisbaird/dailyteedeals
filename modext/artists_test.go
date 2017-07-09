package modext_test

import (
	"testing"

	"github.com/harrisbaird/dailyteedeals/models"
	. "github.com/harrisbaird/dailyteedeals/modext"
	"github.com/nbio/st"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/types"
)

func TestFindOrCreateArtist(t *testing.T) {
	testCases := []struct {
		name               string
		createNew          bool
		artistsCountChange int
	}{
		{"New", true, 1},
		{"Existing", false, 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			RunInTestTransaction(func(db boil.Executor) {
				var artist *models.Artist
				var err error

				if !tt.createNew {
					CreateArtistFixtures(db)
				}

				countDiff := TableCountDiff(db, "artists", func() {
					artist, err = FindOrCreateArtist(db, "theduc", []string{"invalid-domain", "https://www.teepublic.com/user/theduc", "https://www.neatoshop.com/artist/Theduc"})
				})

				st.Expect(t, err, nil)
				st.Reject(t, artist, nil)
				st.Expect(t, countDiff, tt.artistsCountChange)
				st.Expect(t, ValidSlug.MatchString(artist.Slug), true)
			})
		})
	}
}

func TestArtistAppendWhitelistedURLS(t *testing.T) {
	testCases := []struct {
		name    string
		initial types.StringArray
		have    types.StringArray
		want    types.StringArray
	}{
		{"Blank",
			types.StringArray{},
			types.StringArray{},
			types.StringArray{}},
		{"Non-Whitelisted domains",
			types.StringArray{},
			types.StringArray{"http://test.com", "http://bad.com", "http://subdomain.redbubble.com"},
			types.StringArray{}},
		{"Whitelisted domains",
			types.StringArray{},
			types.StringArray{"http://google.com", "http://facebook.com", "http://www.redbubble.com"},
			types.StringArray{"http://facebook.com", "http://redbubble.com"}},
		{"Existing non-whitelist domains are kept",
			types.StringArray{"http://facebook.com", "http://microsoft.com"},
			types.StringArray{"http://google.com", "http://twitter.com"},
			types.StringArray{"http://facebook.com", "http://microsoft.com", "http://twitter.com"}},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			artist := models.Artist{Urls: tt.initial}
			ArtistAppendWhitelistedUrls(&artist, tt.have)
			st.Expect(t, artist.Urls, tt.want)
		})
	}
}

func TestArtistHooks(t *testing.T) {
	RunInTestTransaction(func(db boil.Executor) {
		artist := models.Artist{Name: "   test artist  ",
			Urls: []string{"http://www.google.com", "", "https://othersite.com"},
			Tags: []string{"other artist alias"}}
		err := artist.Insert(db)

		st.Expect(t, err, nil)
		st.Expect(t, ValidSlug.MatchString(artist.Slug), true)
	})
}
