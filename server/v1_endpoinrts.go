package server

import (
	"net/http"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/labstack/echo"
)

func V1ProductsEndpoint(db orm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, productsErr := models.ActiveDeals(db)
		if productsErr != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "deals not found"})
		}

		return c.JSON(http.StatusOK, buildV1Api(products))
	}
}
