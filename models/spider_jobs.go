package models

import (
	"time"

	"github.com/go-pg/pg/orm"
)

type SpiderJob struct {
	ID           int
	SiteID       int
	ScrapydJobID string
	JobType      string
	CreatedAt    time.Time

	Site *Site
}

func FindSpiderJob(db orm.DB, id int) (*SpiderJob, error) {
	var job SpiderJob
	err := db.Model(&job).Where("id=?", id).First()
	return &job, err
}

func CreateSpiderJob(db orm.DB, siteID int, scrapydJobID, jobType string) (*SpiderJob, error) {
	job := SpiderJob{SiteID: siteID, ScrapydJobID: scrapydJobID, JobType: jobType}
	err := db.Insert(&job)
	return &job, err
}
