package models

import (
	"encoding/json"
	"log"
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

func FindSpiderItem(db orm.DB, id int) (*SpiderItem, error) {
	var item SpiderItem
	err := db.Model(&item).Column("spider_item.*", "SpiderJob").Where("spider_item.id=?", id).First()
	return &item, err
}

func CreateSpiderItem(db orm.DB, spiderJobID int, data string) (*SpiderItem, error) {
	item := SpiderItem{SpiderJobID: spiderJobID, ItemData: data}
	err := db.Insert(&item)
	return &item, err
}

func (item *SpiderItem) UpdateError(db orm.DB, err error) {
	item.Error = err.Error()
	if err := db.Update(item); err != nil {
		log.Println(err.Error())
	}
}

func (item *SpiderItem) ParseItemData(db orm.DB, minioConn *utils.MinioConnection) error {
	var data ScrapydItem
	if err := json.Unmarshal([]byte(item.ItemData), &data); err != nil {
		return err
	}

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
		Prices:     data.StringPrices(),
		Active:     data.Active,
		Deal:       data.Deal,
		LastChance: data.LastChance,
		Tags:       data.Tags,
	}

	_, err = db.Model(&product).OnConflict("(design_id, site_id) DO UPDATE").
		Set("url = ?url, prices = ?prices, active = ?active, deal = ?deal, last_chance = ?last_chance").
		Insert()
	if err != nil {
		log.Println(err)
		return err
	}

	db.Model(item).Set("product_id=?", product.ID).Update()

	return product.UpdateImageIfExpired(db, minioConn, data.ImageURL)
}
