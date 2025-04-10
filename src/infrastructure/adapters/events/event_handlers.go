package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"ApiSmart/src/core/application"
	"ApiSmart/src/core/domain/events"
	"ApiSmart/src/core/domain/models"
)

// SensorDataHandler maneja eventos relacionados con datos de sensores
type SensorDataHandler struct {
	sensorService application.SensorService
}

// NewSensorDataHandler crea un nuevo manejador de datos de sensores
func NewSensorDataHandler(sensorService application.SensorService) *SensorDataHandler {
	return &SensorDataHandler{
		sensorService: sensorService,
	}
}

// Handle implementa application.EventHandler
func (h *SensorDataHandler) Handle(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.EventTypeSensorDataCreated:
		return h.handleSensorDataCreated(ctx, event)
	default:
		return fmt.Errorf("tipo de evento no manejado: %s", event.Type)
	}
}

// handleSensorDataCreated procesa eventos de creación de datos de sensores
func (h *SensorDataHandler) handleSensorDataCreated(ctx context.Context, event events.Event) error {
	// Convertir datos del evento a SensorData
	var sensorData models.SensorData
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("error serializando datos del evento: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &sensorData); err != nil {
		return fmt.Errorf("error deserializando a SensorData: %v", err)
	}

	log.Printf("Procesando datos de sensor: ID=%d, Temperatura=%.2f",
		sensorData.ID, sensorData.TemperaturaDHT)

	return nil
}

// AlertHandler maneja eventos relacionados con alertas
type AlertHandler struct {
	sensorService application.SensorService
}

// NewAlertHandler crea un nuevo manejador de alertas
func NewAlertHandler(sensorService application.SensorService) *AlertHandler {
	return &AlertHandler{
		sensorService: sensorService,
	}
}

// Handle implementa application.EventHandler
func (h *AlertHandler) Handle(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.EventTypeSensorThresholdAlert:
		return h.handleSensorAlert(ctx, event)
	default:
		return fmt.Errorf("tipo de evento no manejado: %s", event.Type)
	}
}

// handleSensorAlert procesa eventos de alertas de sensores
func (h *AlertHandler) handleSensorAlert(ctx context.Context, event events.Event) error {
	// Convertir datos del evento a Alert
	var alert models.Alert
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("error serializando datos del evento: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &alert); err != nil {
		return fmt.Errorf("error deserializando a Alert: %v", err)
	}

	// Aquí se implementaría la lógica de notificación (email, SMS, etc.)
	log.Printf("Procesando alerta: Tipo=%s, Mensaje=%s",
		alert.SensorType, alert.Message)

	return nil
}

// UserEventHandler maneja eventos relacionados con usuarios
type UserEventHandler struct {
	authService application.AuthService
}

// NewUserEventHandler crea un nuevo manejador de eventos de usuario
func NewUserEventHandler(authService application.AuthService) *UserEventHandler {
	return &UserEventHandler{
		authService: authService,
	}
}

// Handle implementa application.EventHandler
func (h *UserEventHandler) Handle(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.EventTypeUserRegistered:
		log.Printf("Usuario registrado: %v", event.Data["email"])
		// Aquí se implementaría lógica adicional (envío de email de bienvenida, etc.)
		return nil
	case events.EventTypeUserAuthenticated:
		log.Printf("Usuario autenticado: %v", event.Data["email"])
		// Aquí se implementaría lógica adicional (registro de login, etc.)
		return nil
	default:
		return fmt.Errorf("tipo de evento no manejado: %s", event.Type)
	}
}
