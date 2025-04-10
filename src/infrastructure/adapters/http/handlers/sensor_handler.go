package handlers

import (
	"net/http"
	"strconv"

	"ApiSmart/src/core/application/use_case"
	"ApiSmart/src/core/domain/models"
	"github.com/gin-gonic/gin"
)

// SensorHandler maneja las solicitudes HTTP relacionadas con sensores
type SensorHandler struct {
	sensorUseCase *use_case.SensorUseCase
}

// NewSensorHandler crea una nueva instancia de SensorHandler
func NewSensorHandler(sensorUseCase *use_case.SensorUseCase) *SensorHandler {
	return &SensorHandler{
		sensorUseCase: sensorUseCase,
	}
}

// CreateSensorData crea nuevos datos de sensores
func (h *SensorHandler) CreateSensorData(c *gin.Context) {
	var data models.SensorData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Guardar datos del sensor y generar alertas si es necesario
	if err := h.sensorUseCase.SaveSensorData(c.Request.Context(), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Datos del sensor guardados correctamente",
		"data":    data,
	})
}

// GetAllSensorData obtiene todos los datos de los sensores
func (h *SensorHandler) GetAllSensorData(c *gin.Context) {
	data, err := h.sensorUseCase.GetAllSensorData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetLatestSensorData obtiene los datos más recientes del sensor
func (h *SensorHandler) GetLatestSensorData(c *gin.Context) {
	data, err := h.sensorUseCase.GetLatestSensorData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetAlerts obtiene las alertas de sensores
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

	alerts, err := h.sensorUseCase.GetAlerts(c.Request.Context(), isRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// MarkAlertAsRead marca una alerta como leída
func (h *SensorHandler) MarkAlertAsRead(c *gin.Context) {
	alertID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de alerta inválido"})
		return
	}

	if err := h.sensorUseCase.MarkAlertAsRead(c.Request.Context(), uint(alertID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alerta marcada como leída"})
}
