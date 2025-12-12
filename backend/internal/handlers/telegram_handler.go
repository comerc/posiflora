package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/posiflora/backend/internal/services"
)

type TelegramHandler struct {
	integrationService *services.TelegramIntegrationService
	statusService      *services.StatusService
}

func NewTelegramHandler(
	integrationService *services.TelegramIntegrationService,
	statusService *services.StatusService,
) *TelegramHandler {
	return &TelegramHandler{
		integrationService: integrationService,
		statusService:      statusService,
	}
}

// Connect обрабатывает POST /shops/:shopId/telegram/connect
func (h *TelegramHandler) Connect(c *gin.Context) {
	shopID, err := strconv.ParseInt(c.Param("shopId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return
	}

	var req services.ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	integration, err := h.integrationService.Connect(c.Request.Context(), shopID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

// GetStatus обрабатывает GET /shops/:shopId/telegram/status
func (h *TelegramHandler) GetStatus(c *gin.Context) {
	shopID, err := strconv.ParseInt(c.Param("shopId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopId"})
		return
	}

	// Проверяем параметр mask=disabled для отключения маскирования
	// Если mask=disabled, то maskDisabled=true означает "маскирование отключено" = показывать полный
	maskDisabled := c.Query("mask") == "disabled"

	status, err := h.statusService.GetStatus(c.Request.Context(), shopID, maskDisabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
