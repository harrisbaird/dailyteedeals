package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"

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

	bs := backend.Start(db)
	defer bs.Stop()

	fmt.Printf("Daily Tee Deals is running in %s mode.\n", config.ModeString())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}
