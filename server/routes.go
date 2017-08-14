package server

import (
	"net/http"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Host struct{ Echo *echo.Echo }
type Hosts map[string]*Host

func SetupRoutes(db orm.DB) Hosts {
	apiRouter := ApiRouter(db, newRouter())
	productRedirectRouter := ProductRedirectRouter(db, newRouter())

	hosts := Hosts{}

	hosts[config.App.DomainAPI] = &Host{apiRouter}
	hosts[config.App.DomainGo] = &Host{productRedirectRouter}
	// hs[config.App.DomainWeb] =

	// Also listen locally using lvh.me
	if !config.IsProduction() {
		hosts["api.lvh.me:8080"] = &Host{apiRouter}
		hosts["go.lvh.me:8080"] = &Host{productRedirectRouter}
	}

	return hosts
}

func newRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func ProductRedirectRouter(db orm.DB, e *echo.Echo) *echo.Echo {
	e.GET("/:slug", func(c echo.Context) error {
		product, err := models.FindProductBySlug(db, c.Param("slug"))
		if err != nil {
			return err
		}
		buyURL, err := product.BuyURL(db)
		if err != nil {
			return err
		}
		c.Redirect(http.StatusFound, buyURL)
		return nil
	})
	return e
}

func ApiRouter(db orm.DB, e *echo.Echo) *echo.Echo {
	e.Use(AuthTokenMiddleware(db))

	v1Group := e.Group("/v1")
	v1Group.GET("/products.json", V1ProductsEndpoint(db))

	v2Group := e.Group("/v2")
	v2Group.GET("/deals", V2DealsEndpoint(db))
	v2Group.GET("/designs/:slug", V2DesignEndpoint(db))
	v2Group.GET("/artists/:slug", V2ArtistEndpoint(db))
	v2Group.GET("/sites", V2SiteIndexEndpoint(db))
	v2Group.GET("/sites/:slug", V2SiteShowEndpoint(db))

	return e
}
