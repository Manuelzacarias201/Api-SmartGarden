package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	_ "time"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
	"ApiSmart/pkg/event"
)

// SensorDataHandler maneja eventos relacionados con datos de sensores
type SensorDataHandler struct {
	sensorService ports.SensorService
}

// NewSensorDataHandler crea un nuevo SensorDataHandler
func NewSensorDataHandler(sensorService ports.SensorService) *SensorDataHandler {
	return &SensorDataHandler{
		sensorService: sensorService,
	}
}

// Handle procesa eventos de datos de sensores
func (h *SensorDataHandler) Handle(ctx context.Context, evt event.Event) error {
	switch evt.Type {
	case event.EventTypeSensorDataCreated:
		// Convertir datos del evento a SensorData
		var sensorData domain.SensorData
		dataBytes, err := json.Marshal(evt.Data)
		if err != nil {
			return fmt.Errorf("error marshaling event data: %v", err)
		}

		if err := json.Unmarshal(dataBytes, &sensorData); err != nil {
			return fmt.Errorf("error unmarshaling to SensorData: %v", err)
		}

		// Guardar los datos del sensor
		if err := h.sensorService.SaveSensorData(ctx, &sensorData); err != nil {
			return fmt.Errorf("error saving sensor data: %v", err)
		}

		log.Printf("Sensor data processed and saved: %v", sensorData.ID)
		return nil
	default:
		return fmt.Errorf("unhandled event type: %s", evt.Type)
	}
}

// AlertHandler maneja eventos relacionados con alertas
type AlertHandler struct {
	sensorService ports.SensorService
}

// NewAlertHandler crea un nuevo AlertHandler
func NewAlertHandler(sensorService ports.SensorService) *AlertHandler {
	return &AlertHandler{
		sensorService: sensorService,
	}
}

// Handle procesa eventos de alertas
func (h *AlertHandler) Handle(ctx context.Context, evt event.Event) error {
	switch evt.Type {
	case event.EventTypeSensorThresholdAlert:
		// Convertir datos del evento a Alert
		var alert domain.Alert
		dataBytes, err := json.Marshal(evt.Data)
		if err != nil {
			return fmt.Errorf("error marshaling event data: %v", err)
		}

		if err := json.Unmarshal(dataBytes, &alert); err != nil {
			return fmt.Errorf("error unmarshaling to Alert: %v", err)
		}

		// Procesar la alerta (podría enviar notificaciones, etc.)
		log.Printf("Alert processed: %s - %s", alert.SensorType, alert.Message)
		return nil
	default:
		return fmt.Errorf("unhandled event type: %s", evt.Type)
	}
}

// UserEventHandler maneja eventos relacionados con usuarios
type UserEventHandler struct {
	authService ports.AuthService
}

// NewUserEventHandler crea un nuevo UserEventHandler
func NewUserEventHandler(authService ports.AuthService) *UserEventHandler {
	return &UserEventHandler{
		authService: authService,
	}
}

// Handle procesa eventos de usuarios
func (h *UserEventHandler) Handle(ctx context.Context, evt event.Event) error {
	switch evt.Type {
	case event.EventTypeUserRegistered:
		log.Printf("User registered event: %v", evt.Data["email"])
		// Aquí podrías implementar lógica adicional después del registro
		return nil
	case event.EventTypeUserAuthenticated:
		log.Printf("User authenticated event: %v", evt.Data["email"])
		// Aquí podrías implementar lógica adicional después de la autenticación
		return nil
	default:
		return fmt.Errorf("unhandled event type: %s", evt.Type)
	}
}
