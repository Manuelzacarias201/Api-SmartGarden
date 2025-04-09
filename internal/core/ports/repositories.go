package ports

import (
	"context"

	"ApiSmart/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
}

type SensorRepository interface {
	SaveSensorData(ctx context.Context, data *domain.SensorData) error
	GetAllSensorData(ctx context.Context) ([]domain.SensorData, error)
	GetLatestSensorData(ctx context.Context) (*domain.SensorData, error)
	SaveAlert(ctx context.Context, alert *domain.Alert) error
	GetAlerts(ctx context.Context, isRead *bool) ([]domain.Alert, error)
	MarkAlertAsRead(ctx context.Context, alertID uint) error
}
