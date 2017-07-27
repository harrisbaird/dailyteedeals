package models

import (
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
)

type Site struct {
	ID           int
	Name         string
	Slug         string
	DomainName   string
	AffiliateURL string
	DealScraper  bool
	FullScraper  bool
	Active       bool
	DisplayOrder int

	Products []*Product
}

var jobTypes []string

type SiteJobType int

func (e SiteJobType) String() string {
	return jobTypes[int(e)]
}

func (e SiteJobType) DatabaseField() string {
	return jobTypes[int(e)] + "_scraper"
}

func siteiota(s string) SiteJobType {
	jobTypes = append(jobTypes, s)
	return SiteJobType(len(jobTypes) - 1)
}

var (
	SiteDealJobType = siteiota("deal")
	SiteFullJobType = siteiota("full")
)

func ActiveSitesWithJobType(db orm.DB, jobType SiteJobType) ([]*Site, error) {
	var sites []*Site
	err := db.Model(&sites).Where("active=?", true).Where(jobType.DatabaseField() + "=true").Select()
	return sites, err
}

func ActiveSites(db orm.DB) ([]*Site, error) {
	var sites []*Site
	err := db.Model(&sites).Where("active=?", true).Select()
	return sites, err
}

func FindSiteBySlug(db orm.DB, slug string, page int) (*Site, error) {
	perPage := config.App.ItemsPerPage

	var site Site
	err := db.Model(&site).
		Column("site.*", "Products", "Products.Design", "Products.Design.Artist", "Products.Site").
		Relation("Products", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("product.active=?", true).Offset(perPage * (page - 1)).Limit(perPage), nil
		}).
		Where("site.slug=?", slug).
		First()
	return &site, err
}
