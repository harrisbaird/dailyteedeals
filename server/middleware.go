package server

import (
	"net/http"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/labstack/echo"
)

func AuthTokenMiddleware(db orm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var token string
			if token = c.QueryParam("key"); token == "" {
				// fallback to key in header
				token = c.Request().Header.Get("X-Api-Key")
			}

			if !models.ValidAPIUser(db, token) {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid API key"})
			}

			return next(c)
		}
	}
}
