package models_ext_test

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/harrisbaird/dailyteedeals/models"
	. "github.com/harrisbaird/dailyteedeals/models_ext"
	"github.com/nbio/st"
	"github.com/vattle/sqlboiler/types"
)

func TestFindOrCreateArtist(t *testing.T) {
	db, mock := newSQLMock()
	defer db.Close()

	testCases := []struct {
		name      string
		createNew bool
	}{
		{"New", true},
		{"Existing", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			rows := sqlmock.NewRows([]string{"id", "name", "slug", "urls", "tags"})

			if tt.createNew {
				mock.ExpectQuery("SELECT \\* FROM \"artists\"").WillReturnRows(rows)
				mock.ExpectQuery("INSERT INTO \"artists\"").
					WithArgs("Theduc", sqlmock.AnyArg(), "{\"http://neatoshop.com/artist/Theduc\",\"http://teepublic.com/user/theduc\"}", "{\"theduc\"}").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			} else {
				mock.ExpectQuery("SELECT \\* FROM \"artists\"").WillReturnRows(rows.AddRow(1, "Theduc", "54501-theduc", "{\"http://neatoshop.com/artist/Theduc\",\"http://teepublic.com/user/theduc\"}", "{\"theduc\"}"))
			}

			artist, err := FindOrCreateArtist(db, "theduc", []string{"invalid-domain", "https://www.teepublic.com/user/theduc", "https://www.neatoshop.com/artist/Theduc"})
			st.Expect(t, err, nil)
			st.Expect(t, validateSlug(artist.Slug), true)
			st.Expect(t, mock.ExpectationsWereMet() != nil, false)
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
	db, mock := newSQLMock()
	defer db.Close()

	artist := models.Artist{Name: "   test artist  ",
		Urls: []string{"http://www.google.com", "", "https://othersite.com"},
		Tags: []string{"other artist alias"}}

	mock.ExpectQuery("INSERT INTO \"artists\"").
		WithArgs("Test Artist", sqlmock.AnyArg(), "{\"http://google.com\",\"http://othersite.com\"}", "{\"otherartistalias\",\"testartist\"}").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := artist.Insert(db)

	st.Expect(t, err, nil)
	st.Expect(t, validateSlug(artist.Slug), true)
	st.Expect(t, mock.ExpectationsWereMet(), nil)
}
