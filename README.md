# Daily Tee Deals
Background jobs and and http server for the now defunct dailyteedeals.com, a daily tee aggregator site, written in [Go](https://golang.org).

## Requirements
* [Daily Tee Deals Scrapers](https://github.com/harrisbaird/dailyteedeals_scrapers)
* Postgresql
* Redis
* A S3 bucket or compatible such as Minio for image storage

## Running with Docker
To quickly run all required services using:
* Create an [AWS](https://aws.amazon.com/) bucket to store images (alternatively use [Minio](https://github.com/minio/minio).)
* Install [Docker Compose](https://docs.docker.com/compose/).
* Download and edit the `docker-compose.yml` file.
* Run `docker-compose up` in the directory containing the `docker-compose.yml` file.

A Postgres database dump is available (`dump.sql.gz`)  
Images are not included as they are too large for GitHub (~100GB)



## Frontends

![Frontends](https://github.com/harrisbaird/dailyteedeals/blob/master/assets/frontends.png?raw=true)
