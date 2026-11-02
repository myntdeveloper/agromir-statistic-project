package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myntdeveloper/agromir_boiler_server/internal/service"
)

type BoilerHandler struct {
	service *service.BoilerService
}

func NewBoilerHandler(svc *service.BoilerService) *BoilerHandler {
	return &BoilerHandler{service: svc}
}

func (h *BoilerHandler) ListWarehouses(c *gin.Context) {
	out, err := h.service.ListWarehouses(c.Request.Context())
	if err != nil {
		writeBoilerError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *BoilerHandler) ListBoilers(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	out, err := h.service.ListBoilers(c.Request.Context(), id)
	if err != nil {
		writeBoilerError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *BoilerHandler) GetBoiler(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	out, err := h.service.GetBoiler(c.Request.Context(), id)
	if err != nil {
		writeBoilerError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *BoilerHandler) GetTrend(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	out, err := h.service.GetTrend(c.Request.Context(), id)
	if err != nil {
		writeBoilerError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func parseIDParam(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

func writeBoilerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrWarehouseNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
	case errors.Is(err, service.ErrBoilerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "boiler not found"})
	default:
		log.Printf("%s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
