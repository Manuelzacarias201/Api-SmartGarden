package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Event representa un evento en el sistema
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventHandler es la interfaz para los manejadores de eventos
type EventHandler interface {
	// Handle procesa un evento
	Handle(ctx context.Context, event Event) error
}

// EventDispatcher es responsable de enviar eventos al broker
type EventDispatcher struct {
	broker EventBroker
}

// NewEventDispatcher crea un nuevo EventDispatcher
func NewEventDispatcher(broker EventBroker) *EventDispatcher {
	return &EventDispatcher{
		broker: broker,
	}
}

// Dispatch envía un evento al broker
func (d *EventDispatcher) Dispatch(ctx context.Context, eventType string, topic string, data map[string]interface{}) error {
	event := Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	return d.broker.Publish(ctx, topic, event)
}

// EventTypes define los tipos de eventos disponibles en el sistema
const (
	EventTypeSensorDataCreated    = "sensor.data.created"
	EventTypeSensorThresholdAlert = "sensor.threshold.alert"
	EventTypeSensorDataRequested  = "sensor.data.requested"
	EventTypeUserRegistered       = "user.registered"
	EventTypeUserAuthenticated    = "user.authenticated"
)

// TopicTypes define los topics disponibles en el sistema
const (
	TopicSensorData   = "sensor.data"
	TopicSensorAlerts = "sensor.alerts"
	TopicUserEvents   = "user.events"
)
