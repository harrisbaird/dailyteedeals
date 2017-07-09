package modext

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/disintegration/imaging"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

const (
	imageURLBase      = "https://cdn.dailyteedeals.com/thumbs/%d/%d.jpg?v=%d"
	updateImagesEvery = 7 * 24 * time.Hour
)

var imageTransforms = []utils.Transform{
	utils.Transform{Format: imaging.PNG, MimeType: "image/png", Path: "original/%d.png"},
	utils.Transform{Format: imaging.JPEG, MimeType: "image/jpeg", Size: 300, Path: "thumbs/%d/300.jpg"},
	utils.Transform{Format: imaging.JPEG, MimeType: "image/jpeg", Size: 1200, Path: "thumbs/%d/1200.jpg"},
}

var minioClient = config.NewMinioClient()

func init() {
	models.AddProductHook(boil.BeforeInsertHook, productSaveHook)
	models.AddProductHook(boil.BeforeUpdateHook, productSaveHook)
	models.AddProductHook(boil.BeforeUpsertHook, productSaveHook)
}

// ActiveDeals queries the database for currently active deals
// and eager loads related models.
// It returns a `model.Product` slice and any error encountered.
func ActiveDeals(db boil.Executor) (models.ProductSlice, error) {
	return models.Products(db,
		qm.Load("Site", "Design", "Design.Artist", "Design.Products"),
		// sites join is required for ordering, even though it will be
		// performed again afterwards by the eager load.
		qm.InnerJoin("sites on products.site_id = sites.id"),
		qm.Where("products.active = ?", true),
		qm.Where("products.deal = ?", true),
		qm.OrderBy("sites.display_order ASC, products.site_id ASC, products.last_chance ASC, products.slug ASC"),
	).All()
}

func FindProductBySlug(db boil.Executor, slug string) (*models.Product, error) {
	return models.Products(db, qm.Where("slug=?", slug)).One()
}

// MarkProductsInactive sets products belonging to a site as inactive.
// It either does all deals or non-deals.
func MarkProductsInactive(db boil.Executor, siteID int, deal bool) error {
	return models.
		Products(db, qm.Where("site_id=? AND deal=?", siteID, deal)).
		UpdateAll(models.M{"active": false})
}

func UpdateImageIfExpired(db boil.Executor, p *models.Product, url string) error {
	if time.Since(p.ImageUpdatedAt) < updateImagesEvery {
		return nil
	}

	return UpdateImage(db, p, url)
}

func UpdateImage(db boil.Executor, p *models.Product, url string) error {
	imageFile, err := utils.HTTPGetToTempfile(url)
	if err != nil {
		return err
	}

	image, err := imaging.Open(imageFile.Name())
	if err != nil {
		return err
	}

	bgColor := utils.FindBackgroundColor(image)
	for _, transform := range imageTransforms {
		thumb := utils.TransformImage(image, transform, bgColor)
		imageData := new(bytes.Buffer)
		imaging.Encode(imageData, thumb, transform.Format)

		_, err := minioClient.PutObject(config.S3Bucket(), fmt.Sprintf(transform.Path, p.ID), imageData, transform.MimeType)
		if err != nil {
			return err
		}
	}

	p.ImageBackground = utils.ColorToHex(bgColor)
	p.ImageUpdatedAt = time.Now()
	return p.Update(db)
}

func ProductGoURL(product *models.Product) string {
	return fmt.Sprintf("https://go.dailyteedeals.com/%s", product.Slug)
}

func ProductBuyURL(product *models.Product) string {
	if product.R.Site.AffiliateURL.String != "" {
		return product.URL
	}

	return fmt.Sprintf(product.R.Site.AffiliateURL.String, url.QueryEscape(product.URL))
}

func ProductImageURL(product *models.Product, size int) string {
	return fmt.Sprintf("https://%s/thumbs/%d/%d.jpg?v=%d", config.App.DomainImages,
		product.ID, size, product.ImageUpdatedAt.Unix())
}

func productSaveHook(db boil.Executor, p *models.Product) error {
	if p.Slug == "" {
		err := updateProductSlug(db, p)
		if err != nil {
			return err
		}
	}
	p.Tags = NormalizeTags(p.Tags)
	return nil
}

func updateProductSlug(db boil.Executor, product *models.Product) error {
	design, err := product.Design(db).One()
	if err != nil {
		return err
	}

	site, err := product.Site(db).One()
	if err != nil {
		return err
	}

	product.Slug = MakeSlug(design.Name + " on " + site.Name)

	return nil
}
