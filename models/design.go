package models

import (
	"strings"

	titlecase "github.com/AlasdairF/Titlecase"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// import (
// 	"database/sql"
// 	"strings"

// 	titlecase "github.com/AlasdairF/Titlecase"
// 	"github.com/harrisbaird/dailyteedeals/config"
// 	"github.com/harrisbaird/dailyteedeals/models"
// 	"github.com/vattle/sqlboiler/boil"
// 	"github.com/vattle/sqlboiler/queries/qm"
// )

type Design struct {
	ID           int
	ArtistID     int
	Name         string
	Slug         string
	Description  string   `sql:",notnull"`
	Tags         []string `pg:",array" sql:",notnull"`
	CategoryTags []string `pg:",array" sql:",notnull"`
	Mature       bool     `sql:",notnull"`

	Artist     *Artist
	Products   []*Product
	Categories []*Category `pg:",many2many:category_designs"`
}

func FindOrCreateDesign(db orm.DB, artistID int, name string) (*Design, error) {
	var design Design
	err := db.Model(&design).
		Where("artist_id=?", artistID).
		Where(TagQuery("tags", []string{name}, NormalizeTags)).
		First()

	if err == pg.ErrNoRows {
		design = Design{ArtistID: artistID, Name: name}
		err = db.Insert(&design)
	}

	return &design, err
}

func FindDesignBySlug(db orm.DB, slug string) (*Design, error) {
	var design Design
	err := db.Model(&design).
		Column("design.*", "Artist", "Products", "Products.Site").
		Where("design.slug=?", slug).
		First()
	return &design, err
}

func (d *Design) BeforeInsert(db orm.DB) error {
	return d.normalize(db)
}

func (d *Design) BeforeUpdate(db orm.DB) error {
	return d.normalize(db)
}

func (d *Design) normalize(exec orm.DB) error {
	name := titlecase.English(d.Name)
	d.Name = strings.TrimSpace(name)
	if d.Name != "" {
		d.Tags = append(d.Tags, d.Name)
	}
	if d.Slug == "" {
		d.Slug = MakeSlug(d.Name)
	}
	d.Tags = NormalizeTags(d.Tags)
	d.CategoryTags = NormalizeTags(d.CategoryTags)
	return nil
}
