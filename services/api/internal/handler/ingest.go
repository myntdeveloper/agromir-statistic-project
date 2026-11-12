package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/myntdeveloper/agromir_boiler_server/internal/model"
	"github.com/myntdeveloper/agromir_boiler_server/internal/service"
)

type IngestHandler struct {
	service *service.IngestService
}

func NewIngestHandler(svc *service.IngestService) *IngestHandler {
	return &IngestHandler{service: svc}
}

func (h *IngestHandler) IngestReadings(c *gin.Context) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		rejectIngest(c, raw, "read body: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	body := cleanPLCBody(raw)

	var req model.IngestRequest
	if err := binding.JSON.BindBody(body, &req); err != nil {
		rejectIngest(c, raw, "invalid json/body: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
			"body":    bodySnippet(body, 500),
		})
		return
	}

	accepted, err := h.service.Ingest(c.Request.Context(), req)
	if err != nil {
		rejectIngest(c, body, fmt.Sprintf("plant=%s: %v", req.PlantID, err))
		writeIngestError(c, err)
		return
	}

	log.Printf("accepted packet from %s (%s): %d records, raw=%d clean=%d bytes",
		req.PlantID, c.ClientIP(), accepted, len(raw), len(body))
	c.JSON(http.StatusOK, model.IngestResponse{Accepted: accepted})
}

// cleanPLCBody extracts the JSON object out of a fixed-size S7 Array of Char buffer:
// strips NULs and S7 String length headers before '{', stops at the first NUL after '{'
// or the last '}', and fixes locale decimal commas between digits (1,35 -> 1.35).
func cleanPLCBody(raw []byte) []byte {
	b := raw
	start := bytes.IndexByte(b, '{')
	if start < 0 {
		return bytes.Trim(b, "\x00 \t\r\n")
	}
	b = b[start:]
	if nul := bytes.IndexByte(b, 0); nul >= 0 {
		b = b[:nul]
	}
	if end := bytes.LastIndexByte(b, '}'); end >= 0 {
		b = b[:end+1]
	}
	b = bytes.TrimSpace(b)
	return fixDecimalCommas(b)
}

// fixDecimalCommas turns 1,35 into 1.35 outside JSON strings so PLC locale formatting works.
func fixDecimalCommas(b []byte) []byte {
	var out []byte
	inString := false
	escaped := false
	for i := 0; i < len(b); i++ {
		c := b[i]
		if inString {
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
			} else if c == '"' {
				inString = false
			}
			continue
		}
		if c == '"' {
			inString = true
			continue
		}
		if c == ',' && i > 0 && i+1 < len(b) && isDigit(b[i-1]) && isDigit(b[i+1]) {
			if out == nil {
				out = append([]byte(nil), b...)
			}
			out[i] = '.'
		}
	}
	if out != nil {
		return out
	}
	return b
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (h *IngestHandler) GetPlantStatus(c *gin.Context) {
	plantID := strings.TrimSpace(c.Param("plant_id"))
	if plantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant_id"})
		return
	}

	status, err := h.service.PlantStatus(c.Request.Context(), plantID)
	if err != nil {
		writeIngestError(c, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

// rejectIngest dumps headers+body into API logs and sets X-Ingest-Error for nginx access log.
func rejectIngest(c *gin.Context, body []byte, reason string) {
	c.Header("X-Ingest-Error", headerSafe(reason, 200))

	log.Printf("========== INGEST REJECTED ==========")
	log.Printf("client=%s method=%s path=%s ua=%q",
		c.ClientIP(), c.Request.Method, c.Request.URL.Path, c.Request.UserAgent())
	log.Printf("reason: %s", reason)
	log.Printf("content-type=%q content-length=%q transfer-encoding=%q",
		c.GetHeader("Content-Type"),
		c.GetHeader("Content-Length"),
		c.GetHeader("Transfer-Encoding"))
	log.Printf("headers:\n%s", formatHeaders(c.Request.Header))
	log.Printf("body (%d bytes):\n%s", len(body), formatBody(body))
	log.Printf("=====================================")
}

func formatHeaders(h http.Header) string {
	if len(h) == 0 {
		return "  <none>"
	}
	names := make([]string, 0, len(h))
	for name := range h {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	for _, name := range names {
		for _, value := range h[name] {
			shown := value
			lower := strings.ToLower(name)
			if lower == "authorization" || lower == "header" || lower == "x-api-key" {
				shown = maskSecret(value)
			}
			fmt.Fprintf(&b, "  %s: %s\n", name, shown)
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func maskSecret(value string) string {
	v := strings.TrimSpace(value)
	if len(v) > 7 && strings.EqualFold(v[:7], "Bearer ") {
		return "Bearer " + maskKey(strings.TrimSpace(v[7:]))
	}
	return maskKey(v)
}

func maskKey(key string) string {
	if len(key) <= 4 {
		return "len=" + strconv.Itoa(len(key))
	}
	return key[:4] + "***, len=" + strconv.Itoa(len(key))
}

func formatBody(body []byte) string {
	if len(body) == 0 {
		return "<empty>"
	}
	const max = 64 << 10
	truncated := body
	suffix := ""
	if len(body) > max {
		truncated = body[:max]
		suffix = fmt.Sprintf("\n... truncated, total %d bytes", len(body))
	}
	if !utf8.Valid(truncated) || !mostlyPrintable(truncated) {
		return hexDump(truncated) + suffix
	}
	return string(truncated) + suffix
}

func mostlyPrintable(body []byte) bool {
	bad := 0
	for _, r := range string(body) {
		if r == utf8.RuneError || (r < 32 && r != '\n' && r != '\r' && r != '\t') {
			bad++
		}
	}
	return bad*10 <= len(body)+1
}

func hexDump(body []byte) string {
	var b strings.Builder
	b.WriteString("<non-utf8 or binary; hex>\n")
	for i := 0; i < len(body); i += 16 {
		end := i + 16
		if end > len(body) {
			end = len(body)
		}
		chunk := body[i:end]
		fmt.Fprintf(&b, "%04x  ", i)
		for j := 0; j < 16; j++ {
			if j < len(chunk) {
				fmt.Fprintf(&b, "%02x ", chunk[j])
			} else {
				b.WriteString("   ")
			}
			if j == 7 {
				b.WriteByte(' ')
			}
		}
		b.WriteString(" |")
		for _, c := range chunk {
			if c >= 32 && c < 127 {
				b.WriteByte(c)
			} else {
				b.WriteByte('.')
			}
		}
		b.WriteString("|\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func bodySnippet(body []byte, max int) string {
	if len(body) == 0 {
		return "<empty>"
	}
	if len(body) > max {
		return string(body[:max]) + "...(" + strconv.Itoa(len(body)) + " bytes)"
	}
	if !utf8.Valid(body) {
		return fmt.Sprintf("<binary %d bytes>", len(body))
	}
	return string(body)
}

func headerSafe(s string, max int) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\r' || !unicode.IsPrint(r) {
			b.WriteByte(' ')
			continue
		}
		b.WriteRune(r)
		if b.Len() >= max {
			break
		}
	}
	return strings.TrimSpace(b.String())
}

func writeIngestError(c *gin.Context, err error) {
	var unknown *service.UnknownBoilersError
	switch {
	case errors.As(err, &unknown):
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "unknown boiler_number",
			"numbers": unknown.Numbers,
		})
	case errors.Is(err, service.ErrPlantNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "plant not found"})
	default:
		log.Printf("%s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
