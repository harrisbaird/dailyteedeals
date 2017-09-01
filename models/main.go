package models

import (
	"log"
	"time"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
)

const maxConnectionAttempts = 30 // seconds

func Connect() *pg.DB {
	db := retriableConnect(config.PostgresConnectionOptions(), 0)
	runMigrations(db)
	insertJobTables(db)
	return db
}

func retriableConnect(pgOptions *pg.Options, attempt int) *pg.DB {
	db := pg.Connect(pgOptions)

	if attempt == 1 {
		log.Println("")
	}

	// Execute a simple query to see if postgres is ready
	_, err := db.Exec("select version()")
	if err == nil {
		return db
	}

	if attempt >= maxConnectionAttempts {
		panic("Took too long to acquire database connection")
	}

	log.Print(err)

	time.Sleep(1 + time.Second)

	// Retry the next attempt
	return retriableConnect(pgOptions, attempt+1)
}

// RunInTestTransaction wraps database queries in a transaction.
func RunInTestTransaction(logging bool, fn func(orm.DB)) {
	db := retriableConnect(config.PostgresTestConnectionOptions(), 0)
	migrations.Run(db, "init")
	runMigrations(db)

	if logging {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			if event.Error != nil {
				log.Fatalf("%s %s %v", time.Since(event.StartTime), query, event.Error)
			} else {
				log.Printf("%s %s", time.Since(event.StartTime), query)
			}
		})
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}
	defer tx.Rollback()

	fn(tx)
}

func runMigrations(db *pg.DB) {
	oldVersion, newVersion, err := migrations.Run(db, "up")
	if err != nil {
		panic(err)
	}

	if oldVersion != newVersion {
		log.Panicf("Migrated database from %d to %d\n", oldVersion, newVersion)
	}
}

func insertJobTables(db *pg.DB) {
	db.Exec(`
		DROP TABLE IF EXISTS spider_jobs;
		DROP TABLE IF EXISTS spider_items;

		CREATE TABLE spider_jobs (
			id serial PRIMARY KEY,
			site_id integer NOT NULL references sites(id),
			scrapyd_job_id text NOT NULL,
			job_type text NOT NULL,
			created_at timestamp without time zone NOT NULL DEFAULT now()
		);

		CREATE TABLE spider_items (
			id serial PRIMARY KEY,
			spider_job_id integer NOT NULL references spider_jobs(id),
			product_id integer references products(id),
			item_data text NOT NULL DEFAULT ''::text,
			error text NOT NULL DEFAULT '',
			created_at timestamp without time zone NOT NULL DEFAULT now()
		);
	`)
}
