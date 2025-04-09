package domain

import "time"

type SensorData struct {
	ID             uint      `json:"id"`
	TemperaturaDHT float64   `json:"temperaturaDHT"`
	Luz            float64   `json:"luz"`
	Humedad        float64   `json:"humedad"`
	Humo           float64   `json:"humo"`
	CreatedAt      time.Time `json:"created_at"`
}

type Alert struct {
	ID         uint      `json:"id"`
	SensorID   uint      `json:"sensor_id"`
	SensorType string    `json:"sensor_type"` // "temperatura", "luz", "humedad", "humo"
	Value      float64   `json:"value"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// Umbrales para las alertas
type AlertThresholds struct {
	TemperaturaMax float64
	TemperaturaMin float64
	LuzMax         float64
	LuzMin         float64
	HumedadMax     float64
	HumedadMin     float64
	HumoMax        float64
}

// Valores predeterminados para los umbrales de alertas
var DefaultAlertThresholds = AlertThresholds{
	TemperaturaMax: 30.0,
	TemperaturaMin: 10.0,
	LuzMax:         80.0,
	LuzMin:         20.0,
	HumedadMax:     80.0,
	HumedadMin:     30.0,
	HumoMax:        50.0,
}
