package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/harrisbaird/dailyteedeals/utils"
)

func SetupRoutes(db orm.DB, hs utils.HostSwitch) {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	apiRouter := ApiRouter(db, gin.Default())
	productRedirectRouter := ProductRedirectRouter(db, gin.Default())

	hs[config.App.DomainAPI] = apiRouter
	hs[config.App.DomainGo] = productRedirectRouter

	// Also listen locally using lvh.me
	if !config.IsProduction() {
		hs["api.lvh.me:8080"] = apiRouter
		hs["go.lvh.me:8080"] = productRedirectRouter
	}
}

func ProductRedirectRouter(db orm.DB, r *gin.Engine) *gin.Engine {
	r.GET("/:slug", func(c *gin.Context) {

		product, err := models.FindProductBySlug(db, c.Param("slug"))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		buyURL, err := product.BuyURL(db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(http.StatusFound, buyURL)
	})
	return r
}

func ApiRouter(db orm.DB, r *gin.Engine) *gin.Engine {
	r.Use(AuthTokenMiddleware(db))

	v1Router := r.Group("/v1")
	{
		v1Router.GET("/products.json", V1ProductsEndpoint(db))
	}

	v2Router := r.Group("/v2")
	{
		v2Router.GET("/deals", V2DealsEndpoint(db))
		v2Router.GET("/designs/:slug", V2DesignEndpoint(db))
		v2Router.GET("/artists/:slug", V2ArtistEndpoint(db))
		v2Router.GET("/sites", V2SiteIndexEndpoint(db))
		v2Router.GET("/sites/:slug", V2SiteShowEndpoint(db))
	}

	return r
}
