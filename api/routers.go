package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harrisbaird/dailyteedeals/api/v1"
	"github.com/harrisbaird/dailyteedeals/api/v2"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/modext"
	"github.com/vattle/sqlboiler/boil"
)

func SetupRoutes(db boil.Executor, hs hostSwitch) {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	apiRouter := ApiRouter(db, gin.Default())
	productRedirectRouter := ProductRedirectRouter(db, gin.Default())

	hs[config.App.DomainAPI] = apiRouter
	hs[config.App.DomainGo] = productRedirectRouter

	if config.IsDevelopment() {
		hs["api.lvh.me:8080"] = apiRouter
		hs["go.lvh.me:8080"] = productRedirectRouter
	}
}

func ProductRedirectRouter(db boil.Executor, r *gin.Engine) *gin.Engine {
	r.GET("/:slug", func(c *gin.Context) {

		product, err := modext.FindProductBySlug(db, c.Param("slug"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Redirect(http.StatusFound, modext.ProductBuyURL(product))
	})
	return r
}

func ApiRouter(db boil.Executor, r *gin.Engine) *gin.Engine {
	r.Use(AuthTokenMiddleware(db))

	v1Router := r.Group("/v1")
	{
		v1Router.GET("/products.json", v1.ProductsEndpoint(db))
	}

	v2Router := r.Group("/v2")
	{
		v2Router.GET("/deals", v2.DealsEndpoint(db))
		v2Router.GET("/designs/:slug", v2.DesignEndpoint(db))
		v2Router.GET("/artists/:slug", v2.ArtistEndpoint(db))
		v2Router.GET("/sites", v2.SiteIndexEndpoint(db))
		v2Router.GET("/sites/:slug", v2.SiteShowEndpoint(db))
	}

	return r
}
