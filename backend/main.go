package backend

import (
	"bufio"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/go-pg/pg"
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
	DB         *pg.DB
	SpiderJob  *models.SpiderJob
	SpiderItem *models.SpiderItem
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
	pool.Middleware((*JobContext).FindSpiderJob)

	// Setup job scheduling
	if config.IsProduction() {
		log.Println("Creating periodic backend jobs")
		pool.PeriodicallyEnqueue("0 0 * * * *", "schedule_deal")
		pool.PeriodicallyEnqueue("0 0 0 * * 1", "schedule_full")
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
	log.Printf("Starting job: %s (%s)\n", job.Name, job.ID)

	err := next()
	if err != nil {
		log.Printf("Job errored: %s (%s) - %v\n", job.Name, job.ID, err)
	}

	return err
}

func (c *JobContext) FindSpiderJob(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["spider_job_id"]; ok {
		if err := c.DB.Model(&c.SpiderJob).Where("id=?", job.ArgInt64("spider_job_id")).First(); err != nil {
			return err
		}
	}

	if _, ok := job.Args["spider_item_id"]; ok {
		if err := c.DB.Model(&c.SpiderItem).Where("id=?", job.ArgInt64("spider_item_id")).First(); err != nil {
			return err
		}
	}

	return next()
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
	finished, err := ScrapydIsFinished(c.SpiderJob.ScrapydJobID)
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
	if err := models.MarkProductsInactive(c.DB, c.SpiderJob.SiteID, c.SpiderJob.JobType == models.SiteDealJobType.String()); err != nil {
		return err
	}

	feed, err := ScrapydDownloadFeed(c.SpiderJob.ScrapydJobID)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(feed)
	for scanner.Scan() {
		spiderItem, err := models.CreateSpiderItem(c.DB, int(c.SpiderJob.ID), scanner.Text())
		if err != nil {
			fmt.Println("CreateSpiderItem: " + err.Error())
		}

		enqueuer.Enqueue("parse_item", work.Q{"spider_item_id": spiderItem.ID})
	}

	return scanner.Err()
}

// ParseItem parses the item data, creating the required database rows.
func (c *JobContext) ParseItem(job *work.Job) error {
	err := c.DB.RunInTransaction(func(tx *pg.Tx) error {
		err := c.SpiderItem.ParseItemData(tx, minioConn)
		if err != nil {
			c.SpiderItem.Error = err.Error()
			tx.Update(&c.SpiderItem)
		}
		return err
	})
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
		scrapydJobID, err := ScrapydSchedule(site.Slug + "_" + jobType.String())
		if err != nil {
			log.Println("ScrapydSchedule: " + err.Error())
			continue
		}

		spiderJob, err := models.CreateSpiderJob(c.DB, site.ID, scrapydJobID, jobType.String())
		if err != nil {
			log.Println("CreateSpiderJob: " + err.Error())
			continue
		}

		_, err = enqueuer.EnqueueIn("wait_for_scraper", recheckRateSeconds,
			work.Q{"spider_job_id": spiderJob.ID})
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
