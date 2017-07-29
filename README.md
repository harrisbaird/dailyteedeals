[![Build Status](https://travis-ci.org/harrisbaird/dailyteedeals.svg?branch=master)](https://travis-ci.org/harrisbaird/dailyteedeals)
[![codecov](https://codecov.io/gh/harrisbaird/dailyteedeals/branch/master/graph/badge.svg)](https://codecov.io/gh/harrisbaird/dailyteedeals)
[![Go Report Card](https://goreportcard.com/badge/github.com/harrisbaird/dailyteedeals)](https://goreportcard.com/report/github.com/harrisbaird/dailyteedeals)
[![Docker Build Statu](https://img.shields.io/docker/build/harrisbaird/dailyteedeals.svg)](https://hub.docker.com/r/harrisbaird/dailyteedeals/)
[![](https://images.microbadger.com/badges/image/harrisbaird/dailyteedeals.svg)](https://microbadger.com/images/harrisbaird/dailyteedeals "Get your own image badge on microbadger.com")

# Daily Tee Deals
Background jobs and and http server for [dailyteedeals.com](https://dailyteedeals.com), a daily tee site, written in [Go](https://golang.org).

## Requirements
* [Daily Tee Deals Scrapers](https://github.com/harrisbaird/dailyteedeals_scrapers)
* Postgresql
* Redis
* A S3 bucket or compatible such as Minio

## Running with Docker
To quickly run all required services using:
* Create an [AWS](https://aws.amazon.com/) bucket to store images (alternatively use [Minio](https://github.com/minio/minio).)
* Install [Docker Compose](https://docs.docker.com/compose/).
* Download and edit the `docker-compose.yml` file.
* Run `docker-compose up` in the directory containing the `docker-compose.yml` file.

