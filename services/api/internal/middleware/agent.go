package middleware

import (
	"crypto/hmac"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AgentBearer(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got := bearerToken(c.GetHeader("Authorization"))
		if expected == "" || got == "" || !hmac.Equal([]byte(got), []byte(expected)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
