package config

// RabbitMQConfig contiene la configuración para RabbitMQ
type RabbitMQConfig struct {
	URL          string
	ExchangeName string
	QueueName    string
}

// LoadRabbitMQConfig carga la configuración de RabbitMQ
func LoadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URL:          getEnv("RABBITMQ_URL", "amqp://manuel:manuel@33.14.65.86:5672/"),
		ExchangeName: getEnv("RABBITMQ_EXCHANGE", "smart_api_exchange"),
		QueueName:    getEnv("RABBITMQ_QUEUE", "smart_api_queue"),
	}
}
