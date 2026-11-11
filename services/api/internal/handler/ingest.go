package handler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("ingest from %s: read body: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	var req model.IngestRequest
	if err := binding.JSON.BindBody(body, &req); err != nil {
		log.Printf("ingest from %s: invalid body: %v; body: %s", c.ClientIP(), err, bodySnippet(body))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	accepted, err := h.service.Ingest(c.Request.Context(), req)
	if err != nil {
		log.Printf("ingest from %s (plant %s): %v", c.ClientIP(), req.PlantID, err)
		writeIngestError(c, err)
		return
	}

	log.Printf("accepted packet from %s: %d records", req.PlantID, accepted)
	c.JSON(http.StatusOK, model.IngestResponse{Accepted: accepted})
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

func bodySnippet(body []byte) string {
	const max = 1000
	if len(body) == 0 {
		return "<empty>"
	}
	if len(body) > max {
		return string(body[:max]) + "...(" + strconv.Itoa(len(body)) + " bytes)"
	}
	return string(body)
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
