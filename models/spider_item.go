package models

import (
	"encoding/json"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/utils"
)

type SpiderItem struct {
	ID          int
	SpiderJobID int
	ProductID   int
	ItemData    string
	Error       string
	CreatedAt   time.Time

	SpiderJob *SpiderJob
	Product   *Product
}

type ScrapydItem struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	URL          string            `json:"url"`
	ArtistName   string            `json:"artist_name"`
	ArtistUrls   []string          `json:"artist_urls"`
	Prices       map[string]string `json:"prices"`
	ImageURL     string            `json:"image_url"`
	Tags         []string          `json:"tags"`
	FabricColors []string          `json:"fabric_colors"`
	Active       bool              `json:"active"`
	Deal         bool              `json:"deal"`
	LastChance   bool              `json:"last_chance"`
	Valid        bool              `json:"valid"`
	ExpiresAt    time.Time         `json:"expires_at"`
}

func CreateSpiderItem(db orm.DB, spiderJobID int, data string) (*SpiderItem, error) {
	item := SpiderItem{SpiderJobID: spiderJobID, ItemData: data}
	err := db.Insert(&item)
	return &item, err
}

func (item *SpiderItem) ParseItemData(db orm.DB, minioConn *utils.MinioConnection) error {
	var data ScrapydItem
	json.Unmarshal([]byte(item.ItemData), &data)

	artist, err := FindOrCreateArtist(db, data.ArtistName, data.ArtistUrls)
	if err != nil {
		return err
	}

	design, err := FindOrCreateDesign(db, artist.ID, data.Name)
	if err != nil {
		return err
	}

	product := Product{
		DesignID:   design.ID,
		SiteID:     item.SpiderJob.SiteID,
		URL:        data.URL,
		Prices:     data.Prices,
		Active:     data.Active,
		Deal:       data.Deal,
		LastChance: data.LastChance,
		Tags:       data.Tags,
	}

	_, err = db.Model(&product).OnConflict("(design_id, site_id) DO UPDATE").
		Set("url = ?url, prices = ?prices, active = ?active, deal = ?deal, last_chance = ?last_chance").
		Insert()
	if err != nil {
		return err
	}

	return product.UpdateImageIfExpired(db, minioConn, data.ImageURL)
}
