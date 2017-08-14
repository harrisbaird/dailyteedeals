package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/harrisbaird/dailyteedeals/backend"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/server"
	"github.com/harrisbaird/dailyteedeals/utils"
)

func main() {
	db := models.Connect()
	defer db.Close() // nolint: errcheck

	if err := utils.UpdateRates(); err != nil {
		panic(err)
	}

	servers := server.Start(db)
	defer servers.StopAll()

	backend.Start(db)
	defer backend.Stop()

	log.Println("App is running")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
