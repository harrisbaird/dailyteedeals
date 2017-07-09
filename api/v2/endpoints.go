package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harrisbaird/dailyteedeals/modext"
	"github.com/vattle/sqlboiler/boil"
)

func DealsEndpoint(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := modext.ActiveDeals(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "deals not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": buildV2Products(products)})
	}
}

func DesignEndpoint(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		design, err := modext.FindDesignBySlug(db, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "design not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"design": buildV2Design(design)})
	}
}

func ArtistEndpoint(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		artist, err := modext.FindArtistBySlug(db, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "artist not found", "string": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"artist": buildV2ArtistWithDesigns(artist)})
	}
}

func SiteIndexEndpoint(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		sites, err := modext.ActiveSites(db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "sites not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sites": buildV2Sites(sites)})
	}
}

func SiteShowEndpoint(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		site, err := modext.FindSiteBySlug(db, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "site not found", "string": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"site": buildV2Site(site)})
	}
}
