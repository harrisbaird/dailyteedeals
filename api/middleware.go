package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harrisbaird/dailyteedeals/modext"
	"github.com/vattle/sqlboiler/boil"
)

func AuthTokenMiddleware(db boil.Executor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		if token = c.Query("key"); token == "" {
			// fallback to key in header
			token = c.Request.Header.Get("X-Api-Key")
		}

		if !modext.ValidAPIUser(db, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
	}
}
