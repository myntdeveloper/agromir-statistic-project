package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myntdeveloper/agromir_boiler_server/internal/config"
)

var conf = config.Load()

func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" || !isValidKey(apiKey) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}

func isValidKey(key string) bool {
	return key == conf.APIKey
}
