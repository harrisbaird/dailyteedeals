package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/utils"
)

const httpGracefulTimeout = 5 * time.Second

var server *http.Server

func Start(db orm.DB) {
	log.Println("Starting http server on: " + config.App.HTTPListenAddr)

	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	hs := make(utils.HostSwitch)
	SetupRoutes(db, hs)
	// web.SetupRoutes(db, hs)

	server = &http.Server{
		Addr:    config.App.HTTPListenAddr,
		Handler: hs,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
}

func Stop() {
	log.Println("Stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), httpGracefulTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
