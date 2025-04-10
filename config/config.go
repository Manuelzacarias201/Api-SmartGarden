package config

import (
	"os"
	"strconv"
	"strings"

	"ApiSmart/src/core/tipo_de_datos"
)

// AppConfig contiene toda la configuración de la aplicación
type AppConfig struct {
	ServerPort string
	Database   tipo_de_datos.DatabaseConfig
	JWT        tipo_de_datos.JWTConfig
	CORS       tipo_de_datos.CorsConfig
	RabbitMQ   tipo_de_datos.RabbitMQConfig
}

// LoadConfig carga la configuración desde variables de entorno o valores por defecto
func LoadConfig() *AppConfig {
	return &AppConfig{
		ServerPort: getEnv("SERVER_PORT", "8000"),
		Database: tipo_de_datos.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "manuel"),
			DBName:   getEnv("DB_NAME", "sensores_db"),
		},
		JWT: tipo_de_datos.JWTConfig{
			Secret:      getEnv("JWT_SECRET", "secret_key_cambiar_en_produccion"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		CORS: tipo_de_datos.CorsConfig{
			AllowedOrigins:   strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ","),
			AllowedMethods:   strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS,PATCH"), ","),
			AllowedHeaders:   strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,X-Requested-With"), ","),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 86400),
		},
		RabbitMQ: tipo_de_datos.RabbitMQConfig{
			URL:          getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			ExchangeName: getEnv("RABBITMQ_EXCHANGE", "smart_api_exchange"),
			QueueName:    getEnv("RABBITMQ_QUEUE", "smart_api_queue"),
		},
	}
}

// Helper para obtener variables de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper para obtener variables de entorno como entero
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// Helper para obtener variables de entorno como booleano
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
