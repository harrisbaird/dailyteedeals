package models_test

import (
	"testing"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/database"
	. "github.com/harrisbaird/dailyteedeals/models"
	"github.com/nbio/st"
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
			database.RunInTestTransaction(false, func(db orm.DB) {
				var artist *Artist
				var err error

				if !tt.createNew {
					ImportArtistFixtures(db)
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

func TestFindArtistBySlug(t *testing.T) {
	testCases := []struct {
		name    string
		slug    string
		wantID  int
		wantErr error
	}{
		{"Found", "55555-theduc", 1, nil},
		{"Missing", "missing", 0, pg.ErrNoRows},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			database.RunInTestTransaction(false, func(db orm.DB) {
				ImportArtistFixtures(db)
				artist, err := FindArtistBySlug(db, tt.slug, 1)
				st.Expect(t, err, tt.wantErr)
				st.Expect(t, artist.ID, tt.wantID)
			})
		})
	}
}

// func TestArtistLoadActiveDesigns(t *testing.T) {
// 	database.RunInTestTransaction(true, func(db orm.DB) {
// 		ImportProductFixtures(db)

// 		artist, err := FindArtistBySlug(db, "55555-thehookshot")
// 		st.Expect(t, err, nil)

// 		designs, err := artist.LoadActiveDesigns(db, 1)
// 		st.Expect(t, err, nil)

// 		st.Expect(t, len(designs), 2)
// 		for _, design := range designs {
// 			st.Refute(t, design.Artist, nil)
// 			st.Refute(t, design.Products, nil)
// 			// st.Expect(t, len(design.Products), 1)
// 		}
// 	})
// }

func TestArtistAppendWhitelistedURLS(t *testing.T) {
	testCases := []struct {
		name    string
		initial []string
		have    []string
		want    []string
	}{
		{"Blank",
			[]string{},
			[]string{},
			[]string{}},
		{"Non-Whitelisted domains",
			[]string{},
			[]string{"http://test.com", "http://bad.com"},
			[]string{}},
		{"Whitelisted domains",
			[]string{},
			[]string{"http://google.com", "http://facebook.com"},
			[]string{"http://facebook.com"}},
		{"Existing non-whitelist domains are kept",
			[]string{"http://facebook.com", "http://microsoft.com"},
			[]string{"http://google.com", "http://twitter.com"},
			[]string{"http://facebook.com", "http://microsoft.com", "http://twitter.com"}},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			artist := Artist{Urls: tt.initial}
			artist.ArtistAppendWhitelistedUrls(tt.have)
			st.Expect(t, artist.Urls, tt.want)
		})
	}
}

func TestArtistHooks(t *testing.T) {
	database.RunInTestTransaction(false, func(db orm.DB) {
		artist := Artist{Name: "   test artist  ",
			Urls: []string{"http://www.google.com", "", "https://othersite.com"},
			Tags: []string{"other artist alias"}}
		_, err := db.Model(&artist).Insert()

		st.Expect(t, err, nil)
		st.Expect(t, ValidSlug.MatchString(artist.Slug), true)
	})
}
