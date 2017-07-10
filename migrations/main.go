package main

import (
	"database/sql"
	"log"

	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/pressly/goose"
)

func main() {
	db, err := sql.Open("postgres", config.PostgresConnectionString())
	if err != nil {
		panic(err)
	}

	if err := goose.Run("up", db, "."); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}
