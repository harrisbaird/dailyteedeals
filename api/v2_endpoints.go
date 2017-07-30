package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
)

func V2DealsEndpoint(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := models.ActiveDeals(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "deals not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": buildV2Products(products)})
	}
}

func V2DesignEndpoint(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		design, err := models.FindDesignBySlug(db, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "design not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"design": buildV2Design(design)})
	}
}

func V2ArtistEndpoint(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		artist, err := models.FindArtistBySlug(db, c.Param("slug"), 1)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "artist not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"artist": buildV2ArtistWithDesigns(artist)})
	}
}

func V2SiteIndexEndpoint(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sites, err := models.ActiveSites(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "sites not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sites": buildV2Sites(sites)})
	}
}

func V2SiteShowEndpoint(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		site, err := models.FindSiteBySlug(db, c.Param("slug"), 1)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "site not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"site": buildV2SiteWithProducts(site)})
	}
}
