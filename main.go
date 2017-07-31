package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/harrisbaird/dailyteedeals/backend"
	"github.com/harrisbaird/dailyteedeals/database"
	"github.com/harrisbaird/dailyteedeals/migrations"
	"github.com/harrisbaird/dailyteedeals/server"
	"github.com/harrisbaird/dailyteedeals/utils"
)

func main() {
	migrations.Run()

	db := database.Connect()
	defer db.Close() // nolint: errcheck

	if err := utils.UpdateRates(); err != nil {
		panic(err)
	}

	server.Start(db)
	defer server.Stop()

	backend.Start(db)
	defer backend.Stop()

	log.Println("App is running")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
