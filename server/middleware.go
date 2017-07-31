package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
)

func AuthTokenMiddleware(db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		if token = c.Query("key"); token == "" {
			// fallback to key in header
			token = c.Request.Header.Get("X-Api-Key")
		}

		if !models.ValidAPIUser(db, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
	}
}
