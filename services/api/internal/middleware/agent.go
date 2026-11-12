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
		var reason string
		switch {
		case expected == "":
			reason = "AGENT_API_KEY is not set on server"
		case got == "":
			reason = "no key (expected headers Authorization|Header|X-API-Key)"
		case !hmac.Equal([]byte(got), []byte(expected)):
			reason = "wrong key in header " + source + " (" + maskKey(got) + ")"
		default:
			c.Next()
			return
		}
		log.Printf("ingest auth: %s from %s; headers present: %v", reason, c.ClientIP(), presentHeaders(c))
		c.Header("X-Ingest-Error", reason)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "details": reason})
	}
}

func presentHeaders(c *gin.Context) []string {
	out := make([]string, 0, len(agentKeyHeaders))
	for _, name := range agentKeyHeaders {
		if strings.TrimSpace(c.GetHeader(name)) != "" {
			out = append(out, name)
		}
	}
	return out
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
