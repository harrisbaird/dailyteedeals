package models_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/go-pg/pg/orm"
	. "github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

func TestMarkProductsInactive(t *testing.T) {
	type args struct {
		siteID int
		deal   bool
	}

	testCases := []struct {
		name    string
		args    args
		wantIDS []int
	}{
		{"deal", args{2, true}, []int{1, 2, 4}},
		{"non-deal", args{2, false}, []int{1, 2, 3}},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			RunInTestTransaction(false, func(db orm.DB) {
				ImportProductFixtures(db)
				err := MarkProductsInactive(db, tt.args.siteID, tt.args.deal)
				st.Expect(t, err, nil)

				var products []*Product
				err = db.Model(&products).
					Column("id").
					Where("active=?", true).
					Order("id ASC").
					Select()

				var ids []int
				for _, product := range products {
					ids = append(ids, product.ID)
				}

				st.Expect(t, err, nil)
				st.Expect(t, ids, tt.wantIDS)
				st.Expect(t, err != nil, false)
			})
		})
	}
}

func TestProductGoURL(t *testing.T) {
	product := Product{Slug: "product-slug"}
	st.Expect(t, product.GoURL(), "https://go.dailyteedeals.com/product-slug")
}

func TestBuyURL(t *testing.T) {
	testCases := []struct {
		name       string
		siteID     int
		wantOutput string
	}{
		{"Non affiliate", 1, "http://site.com/product"},
		{"affiliate", 2, "http://affiliatesite.com?url=http%3A%2F%2Fsite.com%2Fproduct"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			RunInTestTransaction(false, func(db orm.DB) {
				ImportSiteFixtures(db)
				product := Product{ID: 1, SiteID: tt.siteID, URL: "http://site.com/product"}
				buyURL, err := product.BuyURL(db)
				st.Expect(t, err, nil)
				st.Expect(t, buyURL, tt.wantOutput)
			})
		})
	}
}

func TestImageURL(t *testing.T) {
	testCases := []struct {
		name       string
		imageType  ImageType
		wantOutput string
	}{
		{"original", OriginalImageType, "https://cdn.dailyteedeals.com/original/1.png?v=1483228800"},
		{"small", SmallImageType, "https://cdn.dailyteedeals.com/thumbs/1/300.jpg?v=1483228800"},
		{"large", LargeImageType, "https://cdn.dailyteedeals.com/thumbs/1/1200.jpg?v=1483228800"},
	}

	product := Product{ID: 1, ImageUpdatedAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			st.Expect(t, product.ImageURL(tt.imageType), tt.wantOutput)
		})
	}
}

func TestProductUpdateImage(t *testing.T) {
	defer gock.Off()
	defer gock.DisableNetworking()
	defer gock.DisableNetworkingFilters()

	fixturePath := utils.ProjectRootPath("models", "testdata", "images", "input.jpg")
	gock.New("http://test.com").Get("/image.jpg").Persist().Reply(200).File(fixturePath)

	// Only allow networking access to minio test server
	gock.EnableNetworking()
	gock.NetworkingFilter(func(req *http.Request) bool {
		return strings.Contains(req.URL.String(), "localhost:9000")
	})

	testCases := []struct {
		name           string
		imageUpdatedAt time.Time
		wantBackground string
		wantUpdate     bool
	}{
		{"valid", time.Now(), "#000000", false},
		{"expired", time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC), "#44261c", true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			utils.RunMinioTest(func(minioConn *utils.MinioConnection) {
				RunInTestTransaction(false, func(db orm.DB) {
					product := ImportProductFixtures(db)[0]
					product.ImageUpdatedAt = tt.imageUpdatedAt

					err := product.UpdateImageIfExpired(db, minioConn, "http://test.com/image.jpg")
					st.Expect(t, err, nil)
					st.Expect(t, product.ImageBackground, tt.wantBackground)

					if tt.wantUpdate {
						TestImageMatchesFixture(t, minioConn, fmt.Sprintf("original/%d.png", product.ID), "original.png")
						TestImageMatchesFixture(t, minioConn, fmt.Sprintf("thumbs/%d/300.jpg", product.ID), "300.jpg")
						TestImageMatchesFixture(t, minioConn, fmt.Sprintf("thumbs/%d/1200.jpg", product.ID), "1200.jpg")
						st.Expect(t, product.ImageBackground, tt.wantBackground)

					} else {
						st.Expect(t, product.ImageUpdatedAt, tt.imageUpdatedAt)
					}
				})
			})
		})
	}
}

func TestProductHooks(t *testing.T) {
	RunInTestTransaction(false, func(db orm.DB) {
		ImportProductFixtures(db)
		product := Product{DesignID: 2, SiteID: 2, URL: "http://test.com"}
		err := db.Insert(&product)
		st.Expect(t, err, nil)
		st.Expect(t, ValidSlug.MatchString(product.Slug), true)
	})
}
