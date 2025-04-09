package handlers

import (
	"net/http"
	"strconv"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	sensorService ports.SensorService
}

func NewSensorHandler(sensorService ports.SensorService) *SensorHandler {
	return &SensorHandler{
		sensorService: sensorService,
	}
}

func (h *SensorHandler) CreateSensorData(c *gin.Context) {
	var data domain.SensorData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Guardar datos del sensor y generar alertas si es necesario
	if err := h.sensorService.SaveSensorData(c.Request.Context(), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Datos del sensor guardados correctamente",
		"data":    data,
	})
}

func (h *SensorHandler) GetAllSensorData(c *gin.Context) {
	data, err := h.sensorService.GetAllSensorData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *SensorHandler) GetLatestSensorData(c *gin.Context) {
	data, err := h.sensorService.GetLatestSensorData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *SensorHandler) GetAlerts(c *gin.Context) {
	// Filtrar alertas por estado (leídas/no leídas)
	var isRead *bool
	isReadParam := c.Query("is_read")
	if isReadParam != "" {
		isReadBool, err := strconv.ParseBool(isReadParam)
		if err == nil {
			isRead = &isReadBool
		}
	}

	alerts, err := h.sensorService.GetAlerts(c.Request.Context(), isRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *SensorHandler) MarkAlertAsRead(c *gin.Context) {
	alertID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de alerta inválido"})
		return
	}

	if err := h.sensorService.MarkAlertAsRead(c.Request.Context(), uint(alertID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alerta marcada como leída"})
}
