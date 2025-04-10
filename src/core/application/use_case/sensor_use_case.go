package use_case

import (
	"context"
	"fmt"
	"log"

	"ApiSmart/src/core/application"
	"ApiSmart/src/core/domain/events"
	"ApiSmart/src/core/domain/models"
)

// SensorUseCase implementa los casos de uso relacionados con sensores
type SensorUseCase struct {
	sensorRepo      application.SensorRepository
	alertService    application.AlertService
	eventDispatcher application.EventDispatcher
}

// NewSensorUseCase crea una nueva instancia de SensorUseCase
func NewSensorUseCase(
	sensorRepo application.SensorRepository,
	alertService application.AlertService,
	eventDispatcher application.EventDispatcher,
) *SensorUseCase {
	return &SensorUseCase{
		sensorRepo:      sensorRepo,
		alertService:    alertService,
		eventDispatcher: eventDispatcher,
	}
}

// SaveSensorData guarda datos de un sensor y genera alertas si es necesario
func (uc *SensorUseCase) SaveSensorData(ctx context.Context, data *models.SensorData) error {
	// Guardar los datos del sensor
	if err := uc.sensorRepo.SaveSensorData(ctx, data); err != nil {
		return err
	}

	// Publicar evento de creación de datos
	if uc.eventDispatcher != nil {
		if err := uc.publishSensorDataCreatedEvent(ctx, data); err != nil {
			log.Printf("Error al publicar evento de creación de datos: %v", err)
		}
	}

	// Verificar si se deben generar alertas
	alerts := uc.alertService.CheckAndCreateAlerts(data)

	// Guardar y publicar las alertas generadas
	for _, alert := range alerts {
		if err := uc.sensorRepo.SaveAlert(ctx, &alert); err != nil {
			return err
		}

		// Publicar evento de alerta
		if uc.eventDispatcher != nil {
			if err := uc.publishAlertEvent(ctx, &alert); err != nil {
				log.Printf("Error al publicar evento de alerta: %v", err)
			}
		}
	}

	return nil
}

// GetAllSensorData obtiene todos los datos de sensores
func (uc *SensorUseCase) GetAllSensorData(ctx context.Context) ([]models.SensorData, error) {
	// Publicar evento de solicitud de datos
	if uc.eventDispatcher != nil {
		uc.publishSensorDataRequestedEvent(ctx, "all")
	}

	return uc.sensorRepo.GetAllSensorData(ctx)
}

// GetLatestSensorData obtiene los datos más recientes del sensor
func (uc *SensorUseCase) GetLatestSensorData(ctx context.Context) (*models.SensorData, error) {
	// Publicar evento de solicitud de datos
	if uc.eventDispatcher != nil {
		uc.publishSensorDataRequestedEvent(ctx, "latest")
	}

	return uc.sensorRepo.GetLatestSensorData(ctx)
}

// GetAlerts obtiene las alertas filtradas por estado
func (uc *SensorUseCase) GetAlerts(ctx context.Context, isRead *bool) ([]models.Alert, error) {
	return uc.sensorRepo.GetAlerts(ctx, isRead)
}

// MarkAlertAsRead marca una alerta como leída
func (uc *SensorUseCase) MarkAlertAsRead(ctx context.Context, alertID uint) error {
	return uc.sensorRepo.MarkAlertAsRead(ctx, alertID)
}

// Métodos auxiliares para publicar eventos

// publishSensorDataCreatedEvent publica un evento de creación de datos del sensor
func (uc *SensorUseCase) publishSensorDataCreatedEvent(ctx context.Context, data *models.SensorData) error {
	eventData := map[string]interface{}{
		"id":             data.ID,
		"temperaturaDHT": data.TemperaturaDHT,
		"luz":            data.Luz,
		"humedad":        data.Humedad,
		"humo":           data.Humo,
		"created_at":     data.CreatedAt,
	}

	return uc.eventDispatcher.Dispatch(
		ctx,
		events.EventTypeSensorDataCreated,
		events.TopicSensorData,
		eventData,
	)
}

// publishAlertEvent publica un evento de alerta
func (uc *SensorUseCase) publishAlertEvent(ctx context.Context, alert *models.Alert) error {
	eventData := map[string]interface{}{
		"id":          alert.ID,
		"sensor_id":   alert.SensorID,
		"sensor_type": alert.SensorType,
		"value":       alert.Value,
		"message":     alert.Message,
		"is_read":     alert.IsRead,
		"created_at":  alert.CreatedAt,
	}

	return uc.eventDispatcher.Dispatch(
		ctx,
		events.EventTypeSensorThresholdAlert,
		events.TopicSensorAlerts,
		eventData,
	)
}

// publishSensorDataRequestedEvent publica un evento de solicitud de datos
func (uc *SensorUseCase) publishSensorDataRequestedEvent(ctx context.Context, requestType string) error {
	eventData := map[string]interface{}{
		"request_type": requestType,
		"timestamp":    fmt.Sprintf("%v", ctx.Value("timestamp")),
	}

	return uc.eventDispatcher.Dispatch(
		ctx,
		events.EventTypeSensorDataRequested,
		events.TopicSensorData,
		eventData,
	)
}
