package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
)

type sensorRepository struct {
	db *sql.DB
}

func NewSensorRepository(db *sql.DB) ports.SensorRepository {
	return &sensorRepository{
		db: db,
	}
}

func (r *sensorRepository) SaveSensorData(ctx context.Context, data *domain.SensorData) error {
	query := `
		INSERT INTO sensor_data (temperatura_dht, luz, humedad, humo, created_at) 
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	data.CreatedAt = now

	result, err := r.db.ExecContext(
		ctx,
		query,
		data.TemperaturaDHT,
		data.Luz,
		data.Humedad,
		data.Humo,
		now,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	data.ID = uint(id)
	return nil
}

func (r *sensorRepository) GetAllSensorData(ctx context.Context) ([]domain.SensorData, error) {
	query := `
		SELECT id, temperatura_dht, luz, humedad, humo, created_at 
		FROM sensor_data 
		ORDER BY created_at DESC 
		LIMIT 1000
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensorDataList []domain.SensorData

	for rows.Next() {
		var data domain.SensorData
		var createdAtStr string

		err := rows.Scan(
			&data.ID,
			&data.TemperaturaDHT,
			&data.Luz,
			&data.Humedad,
			&data.Humo,
			&createdAtStr,
		)

		if err != nil {
			return nil, err
		}

		data.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		sensorDataList = append(sensorDataList, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sensorDataList, nil
}

func (r *sensorRepository) GetLatestSensorData(ctx context.Context) (*domain.SensorData, error) {
	query := `
		SELECT id, temperatura_dht, luz, humedad, humo, created_at 
		FROM sensor_data 
		ORDER BY created_at DESC 
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query)

	var data domain.SensorData
	var createdAtStr string

	err := row.Scan(
		&data.ID,
		&data.TemperaturaDHT,
		&data.Luz,
		&data.Humedad,
		&data.Humo,
		&createdAtStr,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no hay datos de sensores disponibles")
		}
		return nil, err
	}

	data.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

	return &data, nil
}

func (r *sensorRepository) SaveAlert(ctx context.Context, alert *domain.Alert) error {
	query := `
		INSERT INTO alerts (sensor_id, sensor_type, value, message, is_read, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	alert.CreatedAt = now

	result, err := r.db.ExecContext(
		ctx,
		query,
		alert.SensorID,
		alert.SensorType,
		alert.Value,
		alert.Message,
		alert.IsRead,
		now,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	alert.ID = uint(id)
	return nil
}

func (r *sensorRepository) GetAlerts(ctx context.Context, isRead *bool) ([]domain.Alert, error) {
	var query string
	var args []interface{}

	// Base query
	query = `
		SELECT id, sensor_id, sensor_type, value, message, is_read, created_at 
		FROM alerts 
		WHERE 1=1
	`

	// Filtrar por estado de lectura si se especifica
	if isRead != nil {
		query += " AND is_read = ?"
		args = append(args, *isRead)
	}

	// Ordenar por fecha de creación, más recientes primero
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert

	for rows.Next() {
		var alert domain.Alert
		var createdAtStr string

		err := rows.Scan(
			&alert.ID,
			&alert.SensorID,
			&alert.SensorType,
			&alert.Value,
			&alert.Message,
			&alert.IsRead,
			&createdAtStr,
		)

		if err != nil {
			return nil, err
		}

		alert.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		alerts = append(alerts, alert)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (r *sensorRepository) MarkAlertAsRead(ctx context.Context, alertID uint) error {
	query := `UPDATE alerts SET is_read = true WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, alertID)
	return err
}
