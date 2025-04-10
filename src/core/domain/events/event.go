package events

import (
	"time"
)

// Event representa un evento en el sistema
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
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
