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
	truncateTable(db, "artists")
	artists := []Artist{
		{Name: "theduc", Slug: "55555-theduc", Urls: []string{"http://teepublic.com/user/theduc", "http://neatoshop.com/artist/Theduc"}},
		{Name: "thehookshot", Slug: "55555-thehookshot", Urls: []string{"http://society6.com/thehookshot"}},
	}
	if err := db.Insert(&artists); err != nil {
		panic(err)
	}
	return artists
}

func ImportDesignFixtures(db orm.DB) []Design {
	ImportArtistFixtures(db)
	truncateTable(db, "designs")
	designs := []Design{
		{ArtistID: 1, Name: "Summer is here", Slug: "55555-summer-is-here"},
		{ArtistID: 1, Name: "Training Corps"},
		{ArtistID: 2, Name: "Thinking With Chickens"},
		{ArtistID: 2, Name: "iGeek"},
		{ArtistID: 2, Name: "Wizards Rule"},
		{ArtistID: 2, Name: "The Legend of HEY!"},
	}
	if err := db.Insert(&designs); err != nil {
		panic(err)
	}
	return designs
}

func ImportProductFixtures(db orm.DB) []Product {
	ImportDesignFixtures(db)
	ImportSiteFixtures(db)
	truncateTable(db, "products")
	prices := map[string]string{"usd": "1200"}
	products := []Product{
		{DesignID: 1, SiteID: 1, URL: "http://test.com", Active: true, Deal: true, Prices: prices},
		{DesignID: 2, SiteID: 1, Slug: "non_affiliate_link", URL: "http://test.com", Active: true, Deal: false, Prices: prices},
		{DesignID: 3, SiteID: 2, Slug: "affiliate_link", URL: "http://test.com", Active: true, Deal: true, Prices: prices},
		{DesignID: 4, SiteID: 2, URL: "http://test.com", Active: true, Deal: false, Prices: prices},
		{DesignID: 5, SiteID: 2, URL: "http://test.com", Active: false, Deal: false, Prices: prices},
		{DesignID: 6, SiteID: 2, URL: "http://test.com", Active: false, Deal: false, Prices: prices},
	}
	if err := db.Insert(&products); err != nil {
		panic(err)
	}
	return products
}

func ImportSiteFixtures(db orm.DB) []Site {
	truncateTable(db, "sites")
	sites := []Site{
		{Name: "Qwertee", Slug: "qwertee", DomainName: "qwertee.com"},
		{Name: "Teefury", Slug: "teefury", DomainName: "teefury.com", AffiliateURL: "http://affiliatesite.com?url=%s"},
	}
	if err := db.Insert(&sites); err != nil {
		panic(err)
	}
	return sites
}

func ImportUserFixtures(db orm.DB) []User {
	truncateTable(db, "users")
	users := []User{
		{Email: "active_api_user@test.com", APIAccess: true, APIToken: "active_api_user", EncryptedPassword: "password"},
		{Email: "inactive_api_user@test.com", APIAccess: false, APIToken: "inactive_api_user", EncryptedPassword: "password"},
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

// truncateTable ensures table is empty and id nextval is reset to 1.
func truncateTable(db orm.DB, table string) {
	_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
}
