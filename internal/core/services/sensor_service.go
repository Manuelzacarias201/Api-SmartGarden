package services

import (
	"context"
	"fmt"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
	"ApiSmart/pkg/event"
)

// EventDrivenSensorService implementa SensorService con arquitectura basada en eventos
type EventDrivenSensorService struct {
	sensorRepo      ports.SensorRepository
	alertService    ports.AlertService
	eventDispatcher *event.EventDispatcher
}

// NewEventDrivenSensorService crea un nuevo servicio de sensores basado en eventos
func NewEventDrivenSensorService(
	sensorRepo ports.SensorRepository,
	alertService ports.AlertService,
	eventDispatcher *event.EventDispatcher,
) ports.SensorService {
	return &EventDrivenSensorService{
		sensorRepo:      sensorRepo,
		alertService:    alertService,
		eventDispatcher: eventDispatcher,
	}
}

// SaveSensorData guarda los datos del sensor y publica eventos relacionados
func (s *EventDrivenSensorService) SaveSensorData(ctx context.Context, data *domain.SensorData) error {
	// Guardar los datos del sensor
	if err := s.sensorRepo.SaveSensorData(ctx, data); err != nil {
		return err
	}

	// Publicar evento de creación de datos del sensor
	if err := s.publishSensorDataCreatedEvent(ctx, data); err != nil {
		return fmt.Errorf("error publishing sensor data created event: %v", err)
	}

	// Verificar si se deben generar alertas
	alerts := s.alertService.CheckAndCreateAlerts(data)

	// Guardar y publicar las alertas generadas
	for _, alert := range alerts {
		if err := s.sensorRepo.SaveAlert(ctx, &alert); err != nil {
			return err
		}

		// Publicar evento de alerta
		if err := s.publishAlertEvent(ctx, &alert); err != nil {
			return fmt.Errorf("error publishing alert event: %v", err)
		}
	}

	return nil
}

// GetAllSensorData obtiene todos los datos de los sensores
func (s *EventDrivenSensorService) GetAllSensorData(ctx context.Context) ([]domain.SensorData, error) {
	// Publicar evento de solicitud de datos
	s.publishSensorDataRequestedEvent(ctx, "all")

	return s.sensorRepo.GetAllSensorData(ctx)
}

// GetLatestSensorData obtiene los datos más recientes del sensor
func (s *EventDrivenSensorService) GetLatestSensorData(ctx context.Context) (*domain.SensorData, error) {
	// Publicar evento de solicitud de datos
	s.publishSensorDataRequestedEvent(ctx, "latest")

	return s.sensorRepo.GetLatestSensorData(ctx)
}

// GetAlerts obtiene las alertas filtradas por estado
func (s *EventDrivenSensorService) GetAlerts(ctx context.Context, isRead *bool) ([]domain.Alert, error) {
	return s.sensorRepo.GetAlerts(ctx, isRead)
}

// MarkAlertAsRead marca una alerta como leída
func (s *EventDrivenSensorService) MarkAlertAsRead(ctx context.Context, alertID uint) error {
	return s.sensorRepo.MarkAlertAsRead(ctx, alertID)
}

// Métodos auxiliares para publicar eventos

// publishSensorDataCreatedEvent publica un evento de creación de datos del sensor
func (s *EventDrivenSensorService) publishSensorDataCreatedEvent(ctx context.Context, data *domain.SensorData) error {
	eventData := map[string]interface{}{
		"id":             data.ID,
		"temperaturaDHT": data.TemperaturaDHT,
		"luz":            data.Luz,
		"humedad":        data.Humedad,
		"humo":           data.Humo,
		"created_at":     data.CreatedAt,
	}

	return s.eventDispatcher.Dispatch(
		ctx,
		event.EventTypeSensorDataCreated,
		event.TopicSensorData,
		eventData,
	)
}

// publishAlertEvent publica un evento de alerta
func (s *EventDrivenSensorService) publishAlertEvent(ctx context.Context, alert *domain.Alert) error {
	eventData := map[string]interface{}{
		"id":          alert.ID,
		"sensor_id":   alert.SensorID,
		"sensor_type": alert.SensorType,
		"value":       alert.Value,
		"message":     alert.Message,
		"is_read":     alert.IsRead,
		"created_at":  alert.CreatedAt,
	}

	return s.eventDispatcher.Dispatch(
		ctx,
		event.EventTypeSensorThresholdAlert,
		event.TopicSensorAlerts,
		eventData,
	)
}

// publishSensorDataRequestedEvent publica un evento de solicitud de datos
func (s *EventDrivenSensorService) publishSensorDataRequestedEvent(ctx context.Context, requestType string) error {
	eventData := map[string]interface{}{
		"request_type": requestType,
		"timestamp":    ctx.Value("timestamp"),
	}

	return s.eventDispatcher.Dispatch(
		ctx,
		event.EventTypeSensorDataRequested,
		event.TopicSensorData,
		eventData,
	)
}
