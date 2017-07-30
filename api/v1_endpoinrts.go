package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
)

func V1ProductsEndpoint(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, productsErr := models.ActiveDeals(db)
		if productsErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "deals not found"})
			return
		}

		c.JSON(http.StatusOK, buildV1Api(products))
	}
}
