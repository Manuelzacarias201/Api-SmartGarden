package tipo_de_datos

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig define la configuración para JWT
type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

// RabbitMQConfig define la configuración para RabbitMQ
type RabbitMQConfig struct {
	URL          string
	ExchangeName string
	QueueName    string
}

// HTTPConfig define la configuración para el servidor HTTP
type HTTPConfig struct {
	Port int
	Cors CorsConfig
}

// CorsConfig define la configuración CORS
type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}
