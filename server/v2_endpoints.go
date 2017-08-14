package server

import (
	"net/http"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/labstack/echo"
)

type X map[string]interface{}

func V2DealsEndpoint(db orm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := models.ActiveDeals(db)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "deals not found", "string": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"products": buildV2Products(products)})
	}
}

func V2DesignEndpoint(db orm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		design, err := models.FindDesignBySlug(db, c.Param("slug"))
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "design not found", "string": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"design": buildV2Design(design)})
	}
}

func V2ArtistEndpoint(db orm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		artist, err := models.FindArtistBySlug(db, c.Param("slug"), 1)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "artist not found", "string": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"artist": buildV2ArtistWithDesigns(artist)})
	}
}

func V2SiteIndexEndpoint(db orm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sites, err := models.ActiveSites(db)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "sites not found", "string": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"sites": buildV2Sites(sites)})
	}
}

func V2SiteShowEndpoint(db orm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		site, err := models.FindSiteBySlug(db, c.Param("slug"), 1)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "site not found", "string": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"site": buildV2SiteWithProducts(site)})
	}
}
