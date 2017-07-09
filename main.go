package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"

	"github.com/harrisbaird/dailyteedeals/api"
	"github.com/harrisbaird/dailyteedeals/backend"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/vattle/sqlboiler/boil"
)

func main() {
	db, err := sql.Open("postgres", config.PostgresConnectionString())
	if err != nil {
		panic(err)
	}

	if config.IsDevelopment() {
		boil.DebugMode = true
	}

	api.Start(db)
	defer api.Stop()

	backend.Start(db)
	defer backend.Stop()

	log.Printf("App is running in %s mode.\n", config.EnvironmentString())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}
