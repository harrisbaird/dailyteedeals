package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/vattle/sqlboiler/boil"
)

var server *http.Server

func Start(db boil.Executor) {
	log.Println("Starting http server on: " + config.App.HTTPListenAddr)

	hs := make(hostSwitch)
	SetupRoutes(db, hs)

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

// Stop gracefully stops http server
func Stop() {
	log.Println("Stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

type hostSwitch map[string]http.Handler

func (hs hostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}
