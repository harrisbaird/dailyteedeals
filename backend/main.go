package backend

import (
	"bufio"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gocraft/work"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/utils"
)

const (
	jobNamespace       = "dailyteedeals"
	recheckRateSeconds = int64(5)
)

var (
	noRetryOptions = work.JobOptions{MaxFails: 1}
	minioConn      = utils.NewMinioConnection()
	enqueuer       *work.Enqueuer
	pool           *work.WorkerPool
)

type JobContext struct {
	DB *pg.DB
}

func Start(db *pg.DB) {
	log.Println("Starting backend")
	redisPool := &redis.Pool{Dial: func() (redis.Conn, error) { return redis.Dial("tcp", config.App.RedisAddr) }}
	pool = work.NewWorkerPool(JobContext{DB: db}, 10, jobNamespace, redisPool)
	enqueuer = work.NewEnqueuer(jobNamespace, redisPool)

	pool.Middleware(func(c *JobContext, job *work.Job, next work.NextMiddlewareFunc) error {
		c.DB = db
		return next()
	})
	pool.Middleware((*JobContext).Log)

	// Setup job scheduling
	if config.IsProduction() {
		log.Println("Creating periodic backend jobs")
		pool.PeriodicallyEnqueue("0 30 * * * *", "schedule_deal")
		// pool.PeriodicallyEnqueue("0 0 0 * * 1", "schedule_full")
		pool.PeriodicallyEnqueue("0 0 0 * * *", "update_exchange_rates")
	}

	pool.JobWithOptions("schedule_deal", noRetryOptions, (*JobContext).ScheduleDeal)
	pool.JobWithOptions("schedule_full", noRetryOptions, (*JobContext).ScheduleFull)
	pool.JobWithOptions("wait_for_scraper", noRetryOptions, (*JobContext).WaitForScraper)
	pool.JobWithOptions("parse_feed", noRetryOptions, (*JobContext).ParseFeed)
	pool.JobWithOptions("parse_item", noRetryOptions, (*JobContext).ParseItem)
	pool.JobWithOptions("update_exchange_rates", noRetryOptions, (*JobContext).UpdateExchangeRates)
	pool.Start()
}

func Stop() {
	log.Println("Stopping backend")
	pool.Stop()
}

func (c *JobContext) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("[%s] starting job\n", job.Name)

	err := next()
	if err != nil {
		log.Printf("[%s] error: %v", job.Name, err)
	} else {
		log.Printf("[%s] finished", job.Name)
	}

	return err
}

// ScheduleDeal schedules jobs for active sites which have
// deal_scraper enabled.
func (c *JobContext) ScheduleDeal(job *work.Job) error {
	return c.scheduleJobs(models.SiteDealJobType)
}

// ScheduleFull schedules jobs for active sites which have
// full_scraper enabled.
func (c *JobContext) ScheduleFull(job *work.Job) error {
	return c.scheduleJobs(models.SiteFullJobType)
}

// WaitForScraper waits until scrapyd reports that the job is finished
// and schedules a 'parse_feed' job, otherwise reschedules a 'wait_for_scraper
// job in 5 seconds.
func (c *JobContext) WaitForScraper(job *work.Job) error {
	spiderJob, err := models.FindSpiderJob(c.DB, int(job.ArgInt64("spider_job_id")))
	if err != nil {
		return err
	}
	finished, err := ScrapydIsFinished(spiderJob.ScrapydJobID)
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
	spiderJob, err := models.FindSpiderJob(c.DB, int(job.ArgInt64("spider_job_id")))
	if err != nil {
		return err
	}
	if err := models.MarkProductsInactive(c.DB, spiderJob.SiteID, spiderJob.JobType == models.SiteDealJobType.String()); err != nil {
		return err
	}

	feed, err := ScrapydDownloadFeed(spiderJob.ScrapydJobID)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(feed)
	for scanner.Scan() {
		spiderItem, err := models.CreateSpiderItem(c.DB, int(spiderJob.ID), scanner.Text())
		if err != nil {
			fmt.Println("CreateSpiderItem: " + err.Error())
		}

		enqueuer.Enqueue("parse_item", work.Q{"spider_item_id": spiderItem.ID})
	}

	return scanner.Err()
}

// ParseItem parses the item data, creating the required database rows.
func (c *JobContext) ParseItem(job *work.Job) error {
	var err error
	spiderItem, err := models.FindSpiderItem(c.DB, int(job.ArgInt64("spider_item_id")))
	if err != nil {
		return err
	}

	err = c.DB.RunInTransaction(func(tx *pg.Tx) error {
		return spiderItem.ParseItemData(tx, minioConn)
	})

	if err != nil {
		spiderItem.UpdateError(c.DB, err)
	}

	return err
}

func (c *JobContext) UpdateExchangeRates(job *work.Job) error {
	return utils.UpdateRates()
}

func (c *JobContext) scheduleJobs(jobType models.SiteJobType) error {
	log.Printf("Starting jobs for %s sites\n", jobType.String())

	sites, err := models.ActiveSitesWithJobType(c.DB, jobType)
	if err != nil {
		return err
	}

	for _, site := range sites {
		err := scheduleJob(c.DB, site, jobType)
		if err != nil {
			log.Printf("Job scheduling failed for site: %s - %s", site.Name, err.Error())
			continue
		}
	}

	return nil
}

func scheduleJob(db orm.DB, site *models.Site, jobType models.SiteJobType) error {
	scrapydJobID, err := ScrapydSchedule(site.Slug + "_" + jobType.String())
	if err != nil {
		return err
	}

	spiderJob, err := models.CreateSpiderJob(db, site.ID, scrapydJobID, jobType.String())
	if err != nil {
		return err
	}

	_, err = enqueuer.EnqueueIn("wait_for_scraper", recheckRateSeconds,
		work.Q{"spider_job_id": spiderJob.ID})

	return err
}
