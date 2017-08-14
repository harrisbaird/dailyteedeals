package server

import (
	"context"
	"log"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const httpGracefulTimeout = 5 * time.Second

type Server struct {
	Echo *echo.Echo
	Addr string
	TLS  bool
}

func (s *Server) Start() {
	var err error

	if s.TLS {
		log.Printf("Listening on %s", config.App.HTTPSListenAddr)

		if config.IsProduction() {
			// Use LetsEncrypt certificate for production.
			s.Echo.AutoTLSManager.Cache = autocert.DirCache("/etc/letsencrypt")
			s.Echo.AutoTLSManager.Email = config.App.LetsEncryptEmail
			err = s.Echo.StartAutoTLS(config.App.HTTPSListenAddr)
		} else {
			// Use a self-signed certificate for development.
			err = s.Echo.StartTLS(config.App.HTTPSListenAddr, "./certs/cert.pem", "./certs/key.pem")
		}
	} else {
		log.Printf("Listening on %s", config.App.HTTPListenAddr)
		err = s.Echo.Start(config.App.HTTPListenAddr)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), httpGracefulTimeout)
	defer cancel()
	if err := s.Echo.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

type Servers struct {
	Pool []*Server
}

func (s *Servers) Add(e *echo.Echo, addr string, tls bool) {
	e.HideBanner = true
	s.Pool = append(s.Pool, &Server{Echo: e, Addr: addr, TLS: tls})
}

func (s *Servers) StartAll() {
	for _, server := range s.Pool {
		go func(server *Server) { server.Start() }(server)
	}
}

func (s *Servers) StopAll() {
	for _, server := range s.Pool {
		go func(server *Server) { server.Stop() }(server)
	}
}

func Start(db orm.DB) *Servers {
	servers := &Servers{}
	servers.Add(newHTTPRedirectServer(), config.App.HTTPListenAddr, false)
	servers.Add(newHTTPSServer(db), config.App.HTTPSListenAddr, true)
	servers.StartAll()
	return servers
}

func newHTTPRedirectServer() *echo.Echo {
	e := echo.New()
	e.Pre(middleware.HTTPSRedirect())
	return e
}

func newHTTPSServer(db orm.DB) *echo.Echo {
	hosts := SetupRoutes(db)

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			return echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return nil
	})
	return e
}
