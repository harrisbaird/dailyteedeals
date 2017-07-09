package backend

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/work"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/modext"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/types"
)

const (
	jobNamespace       = "dailyteedeals"
	recheckRateSeconds = int64(5)
)

var (
	noRetryOptions = work.JobOptions{MaxFails: 1}
	enqueuer       *work.Enqueuer
	pool           *work.WorkerPool
)

type JobContext struct {
	DB      boil.Executor
	SiteID  int
	JobType modext.SiteJobType
}

func Start(db boil.Executor) {
	log.Println("Starting backend")
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) { return redis.Dial("tcp", config.RedisConnectionString()) }}
	pool = work.NewWorkerPool(JobContext{DB: db}, 10, jobNamespace, redisPool)
	enqueuer = work.NewEnqueuer(jobNamespace, redisPool)

	pool.Middleware(func(c *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
		c.DB = db
		return next()
	})
	pool.Middleware((*JobContext).Log)

	// Setup job scheduling
	if config.IsProduction() {
		pool.PeriodicallyEnqueue("0 40 * * * *", "schedule_deal")
		pool.PeriodicallyEnqueue("0 0 0 * * 1", "schedule_full")
	}

	pool.JobWithOptions("schedule_deal", noRetryOptions, (*JobContext).ScheduleDeal)
	pool.JobWithOptions("schedule_full", noRetryOptions, (*JobContext).ScheduleFull)
	pool.JobWithOptions("wait_for_scraper", noRetryOptions, (*JobContext).WaitForScraper)
	pool.JobWithOptions("parse_feed", noRetryOptions, (*JobContext).ParseFeed)
	pool.JobWithOptions("parse_item", noRetryOptions, (*JobContext).ParseItem)
	pool.Start()
}

func Stop() {
	log.Println("Stopping backend")
	pool.Stop()
}

func (c *JobContext) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Starting job: %s (%s)\n", job.Name, job.ID)

	err := next()
	if err != nil {
		log.Printf("Job errored: %s (%s) - %v\n", job.Name, job.ID, err)
	}

	return err
}

// ScheduleDeal schedules jobs for active sites which have
// deal_scraper enabled.
func (c *JobContext) ScheduleDeal(job *work.Job) error {
	return c.scheduleJobs(modext.SiteDealJobType)
}

// ScheduleFull schedules jobs for active sites which have
// full_scraper enabled.
func (c *JobContext) ScheduleFull(job *work.Job) error {
	return c.scheduleJobs(modext.SiteFullJobType)
}

// WaitForScraper waits until scrapyd reports that the job is finished
// and schedules a 'parse_feed' job, otherwise reschedules a 'wait_for_scraper
// job in 5 seconds.
func (c *JobContext) WaitForScraper(job *work.Job) error {
	scrapydJobID := job.ArgString("scrapyd_job_id")
	finished, err := ScrapydIsFinished(scrapydJobID)
	if err != nil {
		return err
	}

	if finished {
		// Scrapy job has finished, parse item feed.
		enqueuer.Enqueue("parse_feed", job.Args)
	} else {
		// Re-enqueue job with existing args.
		enqueuer.EnqueueUniqueIn("wait_for_scraper", recheckRateSeconds, job.Args)
	}

	return nil
}

// ParseFeed parses downloads and parses the item feed
// and creates a 'parse_item' job for each item.
func (c *JobContext) ParseFeed(job *work.Job) error {
	scrapydJobID := job.ArgString("scrapyd_job_id")
	siteID := job.ArgInt64("site_id")
	deal := job.ArgBool("deal")

	if err := modext.MarkProductsInactive(c.DB, int(siteID), deal); err != nil {
		return err
	}

	feed, err := ScrapydDownloadFeed(scrapydJobID)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(feed)
	for scanner.Scan() {
		enqueuer.Enqueue("parse_item", work.Q{"site_id": siteID, "data": scanner.Text()})
	}

	return scanner.Err()
}

// ParseItem parses a single item, creating the required database rows
// and creating
func (c *JobContext) ParseItem(job *work.Job) error {
	siteID := job.ArgInt64("site_id")
	itemData := job.ArgString("data")

	data := ScrapydItem{}

	if err := json.Unmarshal([]byte(itemData), &data); err != nil {
		return err
	}

	artist, err := modext.FindOrCreateArtist(c.DB, data.ArtistName, data.ArtistUrls)
	if err != nil {
		return err
	}

	design, err := modext.FindOrCreateDesign(c.DB, artist.ID, data.Name)
	if err != nil {
		return err
	}

	// convert prices to sqlboiler hstore format
	prices := make(types.HStore)
	for currency, price := range data.Prices {
		prices[strings.ToUpper(currency)] = sql.NullString{String: fmt.Sprintf("%d", price)}
	}

	product := models.Product{
		DesignID:   design.ID,
		SiteID:     int(siteID),
		URL:        data.URL,
		Prices:     prices,
		Active:     data.Active,
		Deal:       data.Deal,
		LastChance: data.LastChance,
		Tags:       data.Tags,
	}

	err = product.Upsert(c.DB, true, []string{"design_id", "site_id"},
		[]string{"url", "prices", "active", "deal", "last_chance"})
	if err != nil {
		return err
	}

	spew.Dump(err)

	return modext.UpdateImageIfExpired(c.DB, &product, data.ImageURL)
}

func (c *JobContext) scheduleJobs(jobType modext.SiteJobType) error {
	sites, err := modext.ActiveSitesWithJobType(c.DB, jobType)
	if err != nil {
		return err
	}

	for _, site := range sites {
		scrapydJobID, err := ScrapydSchedule(site.Slug + "_" + jobType.String())
		if err != nil {
			log.Printf("Site %d returned error %s", site.ID, err.Error())
		}
		_, err = enqueuer.EnqueueIn("wait_for_scraper", recheckRateSeconds, work.Q{
			"site_id":        site.ID,
			"scrapyd_job_id": scrapydJobID,
			"deal":           jobType == modext.SiteDealJobType,
		})
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
