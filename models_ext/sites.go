package models_ext

import (
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

var jobTypes []string

type SiteJobType int

func (e SiteJobType) String() string {
	return jobTypes[int(e)]
}

func (e SiteJobType) DatabaseField() string {
	return jobTypes[int(e)] + "_scraper"
}

func ciota(s string) SiteJobType {
	jobTypes = append(jobTypes, s)
	return SiteJobType(len(jobTypes) - 1)
}

var (
	SiteDealJobType = ciota("deal")
	SiteFullJobType = ciota("full")
)

func ActiveSitesWithJobType(db boil.Executor, jobType SiteJobType) (models.SiteSlice, error) {
	return models.Sites(db, qm.Where("active=? AND ?=?", true, jobType.DatabaseField(), true)).All()
}
