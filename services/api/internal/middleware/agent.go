package middleware

import (
	"crypto/hmac"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PLC/WinCC senders cannot always set Authorization, so the key is also read from
// "Header" and "X-API-Key", with or without the "Bearer " prefix.
var agentKeyHeaders = []string{"Authorization", "Header", "X-API-Key"}

func AgentBearer(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got, source := agentKey(c)
		switch {
		case expected == "":
			log.Printf("ingest auth: AGENT_API_KEY is not set on server, rejecting %s", c.ClientIP())
		case got == "":
			log.Printf("ingest auth: no key from %s (expected one of headers %v)", c.ClientIP(), agentKeyHeaders)
		case !hmac.Equal([]byte(got), []byte(expected)):
			log.Printf("ingest auth: wrong key from %s in header %q (%s)", c.ClientIP(), source, maskKey(got))
		default:
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}

func agentKey(c *gin.Context) (key, source string) {
	for _, name := range agentKeyHeaders {
		value := strings.TrimSpace(c.GetHeader(name))
		if value == "" {
			continue
		}
		if len(value) > 7 && strings.EqualFold(value[:7], "Bearer ") {
			value = strings.TrimSpace(value[7:])
		}
		return value, name
	}
	return "", ""
}

func maskKey(key string) string {
	if len(key) <= 4 {
		return "len=" + strconv.Itoa(len(key))
	}
	return key[:4] + "***, len=" + strconv.Itoa(len(key))
}
