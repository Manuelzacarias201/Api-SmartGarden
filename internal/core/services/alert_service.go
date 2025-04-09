package services

import (
	"fmt"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
)

type alertService struct {
	thresholds domain.AlertThresholds
}

func NewAlertService() ports.AlertService {
	return &alertService{
		thresholds: domain.DefaultAlertThresholds,
	}
}

func (s *alertService) CheckAndCreateAlerts(data *domain.SensorData) []domain.Alert {
	alerts := []domain.Alert{}

	// Verificar temperatura
	if data.TemperaturaDHT > s.thresholds.TemperaturaMax {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "temperatura",
			Value:      data.TemperaturaDHT,
			Message:    fmt.Sprintf("Temperatura alta: %.2f°C - Ha superado el umbral de %.2f°C", data.TemperaturaDHT, s.thresholds.TemperaturaMax),
			IsRead:     false,
		})
	} else if data.TemperaturaDHT < s.thresholds.TemperaturaMin {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "temperatura",
			Value:      data.TemperaturaDHT,
			Message:    fmt.Sprintf("Temperatura baja: %.2f°C - Por debajo del umbral de %.2f°C", data.TemperaturaDHT, s.thresholds.TemperaturaMin),
			IsRead:     false,
		})
	}

	// Verificar luz
	if data.Luz > s.thresholds.LuzMax {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "luz",
			Value:      data.Luz,
			Message:    fmt.Sprintf("Nivel de luz alto: %.2f%% - Ha superado el umbral de %.2f%%", data.Luz, s.thresholds.LuzMax),
			IsRead:     false,
		})
	} else if data.Luz < s.thresholds.LuzMin {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "luz",
			Value:      data.Luz,
			Message:    fmt.Sprintf("Nivel de luz bajo: %.2f%% - Por debajo del umbral de %.2f%%", data.Luz, s.thresholds.LuzMin),
			IsRead:     false,
		})
	}

	// Verificar humedad
	if data.Humedad > s.thresholds.HumedadMax {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "humedad",
			Value:      data.Humedad,
			Message:    fmt.Sprintf("Nivel de humedad alto: %.2f%% - Ha superado el umbral de %.2f%%", data.Humedad, s.thresholds.HumedadMax),
			IsRead:     false,
		})
	} else if data.Humedad < s.thresholds.HumedadMin {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "humedad",
			Value:      data.Humedad,
			Message:    fmt.Sprintf("Nivel de humedad bajo: %.2f%% - Por debajo del umbral de %.2f%%", data.Humedad, s.thresholds.HumedadMin),
			IsRead:     false,
		})
	}

	// Verificar humo
	if data.Humo > s.thresholds.HumoMax {
		alerts = append(alerts, domain.Alert{
			SensorID:   data.ID,
			SensorType: "humo",
			Value:      data.Humo,
			Message:    fmt.Sprintf("Nivel de humo alto: %.2f%% - Ha superado el umbral de %.2f%%", data.Humo, s.thresholds.HumoMax),
			IsRead:     false,
		})
	}

	return alerts
}
