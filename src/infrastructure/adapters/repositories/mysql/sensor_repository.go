package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ApiSmart/src/core/application"
	"ApiSmart/src/core/domain/models"
)

// SensorRepository implementa application.SensorRepository
type SensorRepository struct {
	db *sql.DB
}

// NewSensorRepository crea una nueva instancia de SensorRepository
func NewSensorRepository(db *sql.DB) application.SensorRepository {
	return &SensorRepository{
		db: db,
	}
}

// SaveSensorData guarda datos de un sensor en la base de datos
func (r *SensorRepository) SaveSensorData(ctx context.Context, data *models.SensorData) error {
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

// GetAllSensorData obtiene todos los datos de sensores
func (r *SensorRepository) GetAllSensorData(ctx context.Context) ([]models.SensorData, error) {
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

	var sensorDataList []models.SensorData

	for rows.Next() {
		var data models.SensorData
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

// GetLatestSensorData obtiene los datos más recientes del sensor
func (r *SensorRepository) GetLatestSensorData(ctx context.Context) (*models.SensorData, error) {
	query := `
		SELECT id, temperatura_dht, luz, humedad, humo, created_at 
		FROM sensor_data 
		ORDER BY created_at DESC 
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query)

	var data models.SensorData
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

// SaveAlert guarda una alerta en la base de datos
func (r *SensorRepository) SaveAlert(ctx context.Context, alert *models.Alert) error {
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

// GetAlerts obtiene las alertas filtradas por estado
func (r *SensorRepository) GetAlerts(ctx context.Context, isRead *bool) ([]models.Alert, error) {
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

	var alerts []models.Alert

	for rows.Next() {
		var alert models.Alert
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

// MarkAlertAsRead marca una alerta como leída
func (r *SensorRepository) MarkAlertAsRead(ctx context.Context, alertID uint) error {
	query := `UPDATE alerts SET is_read = true WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, alertID)
	return err
}
