package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/api"
	"github.com/harrisbaird/dailyteedeals/backend"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/database"
	"github.com/harrisbaird/dailyteedeals/migrations"
	"github.com/harrisbaird/dailyteedeals/utils"
)

const httpGracefulTimeout = 5 * time.Second

func main() {
	migrations.Run()

	db := database.Connect()
	defer db.Close() // nolint: errcheck

	if err := utils.UpdateRates(); err != nil {
		panic(err)
	}

	startHTTPServer(db)
	defer gracefullyStopHTTPServer()

	backend.Start(db)
	defer backend.Stop()

	log.Println("App is running")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

var server *http.Server

func startHTTPServer(db orm.DB) {
	log.Println("Starting http server on: " + config.App.HTTPListenAddr)

	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	hs := make(utils.HostSwitch)
	api.SetupRoutes(db, hs)
	// web.SetupRoutes(db, hs)

	server = &http.Server{
		Addr:    config.App.HTTPListenAddr,
		Handler: hs,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
}

func gracefullyStopHTTPServer() {
	log.Println("Stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), httpGracefulTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
