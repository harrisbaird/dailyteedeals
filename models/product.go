package models

import (
	"bytes"
	"net/url"
	"time"

	"fmt"

	"github.com/disintegration/imaging"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/utils"
)

var imageTypes []utils.Transform

type ImageType int

func (t ImageType) Transform() utils.Transform {
	return imageTypes[int(t)]
}

func imageIota(t utils.Transform) ImageType {
	imageTypes = append(imageTypes, t)
	return ImageType(len(imageTypes) - 1)
}

var (
	OriginalImageType = imageIota(utils.Transform{Format: imaging.PNG, MimeType: "image/png", Path: "original/%d.png"})
	SmallImageType    = imageIota(utils.Transform{Format: imaging.JPEG, MimeType: "image/jpeg", Size: 300, Path: "thumbs/%d/300.jpg"})
	LargeImageType    = imageIota(utils.Transform{Format: imaging.JPEG, MimeType: "image/jpeg", Size: 1200, Path: "thumbs/%d/1200.jpg"})
)

type Product struct {
	ID              int
	DesignID        int
	SiteID          int
	Slug            string `sql:",notnull"`
	URL             string
	Active          bool                              `sql:",notnull"`
	Deal            bool                              `sql:",notnull"`
	LastChance      bool                              `sql:",notnull"`
	Tags            []string                          `pg:",array" sql:"tags,DEFAULT:'{}',notnull"`
	Prices          map[string]string                 `pg:",hstore"`
	ConvertedPrices map[string]utils.ApproximatePrice `sql:"-"`
	ImageBackground string
	ImageUpdatedAt  time.Time
	ExpiresAt       pg.NullTime
	ActiveAt        pg.NullTime

	Design *Design
	Site   *Site
}

const updateImagesEvery = 7 * 24 * time.Hour

// ActiveDeals queries the database for currently active deals
// and eager loads related models.
// It returns a `model.Product` slice and any error encountered.
func ActiveDeals(db orm.DB) ([]*Product, error) {
	var products []*Product
	err := db.Model(&products).
		Column("product.*", "Site", "Design", "Design.Artist").
		Where("product.active=? AND product.deal=?", true, true).
		Order("site.display_order ASC", "product.site_id ASC", "product.last_chance ASC", "product.slug ASC").
		Select()
	return products, err
}

func FindProductBySlug(db orm.DB, slug string) (*Product, error) {
	var product Product
	err := db.Model(&product).
		Column("product.*", "Site", "Design", "Design.Artist").
		Where("product.slug=?", slug).
		First()
	return &product, err
}

// MarkProductsInactive sets products belonging to a site as inactive.
// It either does all deals or non-deals.
func MarkProductsInactive(db orm.DB, siteID int, deal bool) error {
	_, err := db.Exec("UPDATE products SET active=? WHERE site_id=? AND deal=?", false, siteID, deal)
	return err
}

func (p *Product) GoURL() string {
	return fmt.Sprintf("https://go.dailyteedeals.com/%s", p.Slug)
}

func (p *Product) BuyURL(db orm.DB) (string, error) {
	var site Site
	err := db.Model(&site).Where("id=?", p.SiteID).First()

	if err != nil || site.AffiliateURL == "" {
		return p.URL, err
	}

	return fmt.Sprintf(site.AffiliateURL, url.QueryEscape(p.URL)), nil
}

func (p *Product) ImageURL(imageType ImageType) string {
	transform := imageType.Transform()
	return fmt.Sprintf("https://%s/"+transform.Path+"?v=%d", config.App.DomainImages,
		p.ID, p.ImageUpdatedAt.Unix())
}

func (p *Product) SmallImageURL() string {
	return p.ImageURL(SmallImageType)
}

func (p *Product) LargeImageURL() string {
	return p.ImageURL(LargeImageType)
}

func (p *Product) UpdateImageIfExpired(db orm.DB, minioConn *utils.MinioConnection, url string) error {
	if time.Since(p.ImageUpdatedAt) < updateImagesEvery {
		return nil
	}

	return p.UpdateImage(db, minioConn, url)
}

func (p *Product) UpdateImage(db orm.DB, minioConn *utils.MinioConnection, url string) error {
	imageFile, err := utils.HTTPGetToTempfile(url)
	if err != nil {
		return err
	}

	image, err := imaging.Open(imageFile.Name())
	if err != nil {
		return err
	}

	bgColor := utils.FindBackgroundColor(image)
	for _, transform := range imageTypes {
		thumb := utils.TransformImage(image, transform, bgColor)
		imageData := new(bytes.Buffer)
		imaging.Encode(imageData, thumb, transform.Format)
		_, err := minioConn.Client.PutObject(minioConn.Bucket, fmt.Sprintf(transform.Path, p.ID), imageData, transform.MimeType)
		if err != nil {
			return err
		}
	}

	p.ImageBackground = utils.ColorToHex(bgColor)
	p.ImageUpdatedAt = time.Now()
	return db.Update(&p)
}

func (p *Product) BeforeInsert(db orm.DB) error {
	return p.normalize(db)
}

func (p *Product) BeforeUpdate(db orm.DB) error {
	return p.normalize(db)
}

func (p *Product) AfterQuery(db orm.DB) error {
	p.ConvertedPrices = utils.ConvertPrices(p.Prices)
	return nil
}

func (p *Product) normalize(db orm.DB) error {
	if p.Slug == "" {
		err := p.updateProductSlug(db)
		if err != nil {
			return err
		}
	}
	p.Tags = NormalizeTags(p.Tags)
	return nil
}

func (p *Product) updateProductSlug(db orm.DB) error {
	design := Design{ID: p.DesignID}
	err := db.Select(&design)
	if err != nil {
		return err
	}

	site := Site{ID: p.SiteID}
	err = db.Select(&site)
	if err != nil {
		return err
	}

	p.Slug = MakeSlug(design.Name + " on " + site.Name)
	return nil
}
