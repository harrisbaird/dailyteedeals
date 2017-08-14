package server

import (
	"context"
	"log"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/labstack/echo"
)

const httpGracefulTimeout = 5 * time.Second

var server *echo.Echo

func Start(db orm.DB) {
	log.Println("Starting http server on: " + config.App.HTTPListenAddr)

	server = newServer(db)

	go func() {
		if err := server.Start(config.App.HTTPListenAddr); err != nil {
			log.Println("HTTP Server", err)
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

func newServer(db orm.DB) *echo.Echo {
	hosts := SetupRoutes(db)

	e := echo.New()
	e.HideBanner = true

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})
	return e
}
