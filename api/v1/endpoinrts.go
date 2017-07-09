package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harrisbaird/dailyteedeals/modext"
	"github.com/vattle/sqlboiler/boil"
)

func ProductsEndpoint(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, productsErr := modext.ActiveDeals(db)
		if productsErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "deals not found"})
			return
		}

		c.JSON(http.StatusOK, buildV1Api(products))
	}
}
