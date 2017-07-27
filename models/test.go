package models

import (
	"image"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

// ValidSlug regexp in "55555-slug-string" format.
var ValidSlug = regexp.MustCompile(`^\d{5}-[a-z0-9-]+`)

func TestImageMatchesFixture(t *testing.T, conn *utils.MinioConnection, objectPath, fixtureName string) {
	tmpfile, err := ioutil.TempFile("", "image")
	st.Assert(t, err, nil)

	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	err = conn.Client.FGetObject(conn.Bucket, objectPath, tmpfile.Name())
	st.Assert(t, err, nil)

	fixturePath := utils.ProjectRootPath("models", "testdata", "images", "output", fixtureName)
	fixture, err := imaging.Open(fixturePath)
	st.Assert(t, err, nil)

	generated, err := imaging.Open(tmpfile.Name())
	st.Assert(t, err, nil)

	// Compare generated image to fixture
	st.Assert(t, compareNRGBA(fixture.(*image.NRGBA), generated.(*image.NRGBA), 0), true)
}

func TableCountDiff(db orm.DB, table string, testFn func()) int {
	countBefore := getTableCount(db, table)
	testFn()
	return getTableCount(db, table) - countBefore
}

func ImportArtistFixtures(db orm.DB) []Artist {
	artists := []Artist{
		{ID: 1, Name: "theduc", Slug: "55555-theduc", Urls: []string{"http://teepublic.com/user/theduc", "http://neatoshop.com/artist/Theduc"}},
		{ID: 2, Name: "thehookshot", Slug: "55555-thehookshot", Urls: []string{"http://society6.com/thehookshot"}},
	}
	if err := db.Insert(&artists); err != nil {
		panic(err)
	}
	return artists
}

func ImportDesignFixtures(db orm.DB) []Design {
	ImportArtistFixtures(db)
	designs := []Design{
		{ID: 1, ArtistID: 1, Name: "Summer is here", Slug: "55555-summer-is-here"},
		{ID: 2, ArtistID: 1, Name: "Training Corps"},
		{ID: 3, ArtistID: 2, Name: "Thinking With Chickens"},
		{ID: 4, ArtistID: 2, Name: "iGeek"},
		{ID: 5, ArtistID: 2, Name: "Wizards Rule"},
		{ID: 6, ArtistID: 2, Name: "The Legend of HEY!"},
	}
	if err := db.Insert(&designs); err != nil {
		panic(err)
	}
	return designs
}

func ImportProductFixtures(db orm.DB) []Product {
	ImportDesignFixtures(db)
	ImportSiteFixtures(db)
	prices := map[string]string{"usd": "1200"}
	products := []Product{
		{ID: 1, DesignID: 1, SiteID: 1, URL: "http://test.com", Active: true, Deal: true, Prices: prices},
		{ID: 2, DesignID: 2, SiteID: 1, Slug: "non_affiliate_link", URL: "http://test.com", Active: true, Deal: false, Prices: prices},
		{ID: 3, DesignID: 3, SiteID: 2, Slug: "affiliate_link", URL: "http://test.com", Active: true, Deal: true, Prices: prices},
		{ID: 4, DesignID: 4, SiteID: 2, URL: "http://test.com", Active: true, Deal: false, Prices: prices},
		{ID: 5, DesignID: 5, SiteID: 2, URL: "http://test.com", Active: false, Deal: false, Prices: prices},
		{ID: 6, DesignID: 6, SiteID: 2, URL: "http://test.com", Active: false, Deal: false, Prices: prices},
	}
	if err := db.Insert(&products); err != nil {
		panic(err)
	}
	return products
}

func ImportSiteFixtures(db orm.DB) []Site {
	sites := []Site{
		{ID: 1, Name: "Qwertee", Slug: "qwertee", DomainName: "qwertee.com"},
		{ID: 2, Name: "Teefury", Slug: "teefury", DomainName: "teefury.com", AffiliateURL: "http://affiliatesite.com?url=%s"},
	}
	if err := db.Insert(&sites); err != nil {
		panic(err)
	}
	return sites
}

func ImportUserFixtures(db orm.DB) []User {
	users := []User{
		{ID: 1, Email: "active_api_user@test.com", APIAccess: true, APIToken: "active_api_user", EncryptedPassword: "password"},
		{ID: 2, Email: "inactive_api_user@test.com", APIAccess: false, APIToken: "inactive_api_user", EncryptedPassword: "password"},
	}

	if err := db.Insert(&users); err != nil {
		panic(err)
	}
	return users
}

func getTableCount(db orm.DB, table string) int {
	count, err := db.Model().Table(table).Count()
	if err != nil {
		panic(err)
	}

	return count
}

// compareNRGBA
// Source: https://github.com/disintegration/imaging
func compareNRGBA(img1, img2 *image.NRGBA, delta int) bool {
	if !img1.Rect.Eq(img2.Rect) {
		return false
	}

	if len(img1.Pix) != len(img2.Pix) {
		return false
	}

	for i := 0; i < len(img1.Pix); i++ {
		if absint(int(img1.Pix[i])-int(img2.Pix[i])) > delta {
			return false
		}
	}

	return true
}

// absint returns the absolute value of i.
func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
