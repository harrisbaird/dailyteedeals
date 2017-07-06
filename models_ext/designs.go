package models_ext

import (
	"database/sql"
	"strings"

	titlecase "github.com/AlasdairF/Titlecase"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

func init() {
	models.AddDesignHook(boil.BeforeInsertHook, designSaveHook)
	models.AddDesignHook(boil.BeforeUpdateHook, designSaveHook)
}

func FindOrCreateDesign(db boil.Executor, artistID int, name string) (*models.Design, error) {
	design, err := models.Designs(db,
		qm.Where("artist_id=?", artistID),
		qm.Where(TagQuery("tags", []string{name}, NormalizeTags)),
	).One()

	if err == sql.ErrNoRows {
		design = &models.Design{ArtistID: artistID, Name: name}
		err = design.Insert(db)
	}

	return design, err
}

func designSaveHook(exec boil.Executor, d *models.Design) error {
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
