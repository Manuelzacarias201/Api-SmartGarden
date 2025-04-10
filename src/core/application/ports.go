package application

import (
	"context"

	"ApiSmart/src/core/domain/events"
	"ApiSmart/src/core/domain/models"
)

// UserRepository define la interfaz para el acceso a datos de usuarios
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}

// SensorRepository define la interfaz para el acceso a datos de sensores
type SensorRepository interface {
	SaveSensorData(ctx context.Context, data *models.SensorData) error
	GetAllSensorData(ctx context.Context) ([]models.SensorData, error)
	GetLatestSensorData(ctx context.Context) (*models.SensorData, error)
	SaveAlert(ctx context.Context, alert *models.Alert) error
	GetAlerts(ctx context.Context, isRead *bool) ([]models.Alert, error)
	MarkAlertAsRead(ctx context.Context, alertID uint) error
}

// AuthService define la interfaz para el servicio de autenticación
type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	ValidateToken(token string) (uint, error)
}

// SensorService define la interfaz para el servicio de sensores
type SensorService interface {
	SaveSensorData(ctx context.Context, data *models.SensorData) error
	GetAllSensorData(ctx context.Context) ([]models.SensorData, error)
	GetLatestSensorData(ctx context.Context) (*models.SensorData, error)
	GetAlerts(ctx context.Context, isRead *bool) ([]models.Alert, error)
	MarkAlertAsRead(ctx context.Context, alertID uint) error
}

// AlertService define la interfaz para el servicio de alertas
type AlertService interface {
	CheckAndCreateAlerts(data *models.SensorData) []models.Alert
}

// EventDispatcher define la interfaz para el despachador de eventos
type EventDispatcher interface {
	Dispatch(ctx context.Context, eventType string, topic string, data map[string]interface{}) error
}

// EventHandler define la interfaz para el manejador de eventos
type EventHandler interface {
	Handle(ctx context.Context, event events.Event) error
}

// EventBroker define la interfaz para el broker de eventos
type EventBroker interface {
	Publish(ctx context.Context, topic string, event events.Event) error
	Subscribe(topic string, handler EventHandler) error
	Close() error
}
