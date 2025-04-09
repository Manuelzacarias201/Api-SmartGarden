package ports

import (
	"context"

	"ApiSmart/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	ValidateToken(token string) (uint, error)
}

type SensorService interface {
	SaveSensorData(ctx context.Context, data *domain.SensorData) error
	GetAllSensorData(ctx context.Context) ([]domain.SensorData, error)
	GetLatestSensorData(ctx context.Context) (*domain.SensorData, error)
	GetAlerts(ctx context.Context, isRead *bool) ([]domain.Alert, error)
	MarkAlertAsRead(ctx context.Context, alertID uint) error
}

type AlertService interface {
	CheckAndCreateAlerts(data *domain.SensorData) []domain.Alert
}

// EventPublisher define la interfaz para componentes que publican eventos
type EventPublisher interface {
	PublishSensorDataCreated(ctx context.Context, data *domain.SensorData) error
	PublishSensorAlert(ctx context.Context, alert *domain.Alert) error
	PublishUserRegistered(ctx context.Context, userID uint, username, email string) error
	PublishUserAuthenticated(ctx context.Context, userID uint, username, email string) error
}

// EventSubscriber define la interfaz para componentes que se suscriben a eventos
type EventSubscriber interface {
	SubscribeToSensorData(handler func(data *domain.SensorData) error) error
	SubscribeToSensorAlerts(handler func(alert *domain.Alert) error) error
	SubscribeToUserEvents(handler func(eventType string, userID uint, username, email string) error) error
}
