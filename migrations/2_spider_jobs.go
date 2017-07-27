package migrations

import "github.com/go-pg/migrations"

func init() {
	migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS spider_jobs (
			id serial PRIMARY KEY,
			site_id integer NOT NULL references sites(id),
			scrapyd_job_id text NOT NULL,
			job_type text NOT NULL,
			created_at timestamp without time zone NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS spider_items (
			id serial PRIMARY KEY,
			spider_job_id integer NOT NULL references spider_jobs(id),
			product_id integer references products(id),
			item_data text NOT NULL DEFAULT ''::text,
			error text NOT NULL DEFAULT '',
			created_at timestamp without time zone NOT NULL DEFAULT now()
		);
	`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
		DROP TABLE spider_jobs;
		DROP TABLE spider_items;
	`)
		return err
	})
}
