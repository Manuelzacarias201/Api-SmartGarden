package event

import (
	"ApiSmart/config"
	"log"
)

// InitEventSystem inicializa el sistema de eventos
func InitEventSystem(cfg config.RabbitMQConfig) (*EventDispatcher, *RabbitMQBroker, error) {
	// Crear broker de RabbitMQ
	rabbitMQConfig := RabbitMQConfig{
		URL:          cfg.URL,
		ExchangeName: cfg.ExchangeName,
		QueueName:    cfg.QueueName,
	}

	broker, err := NewRabbitMQBroker(rabbitMQConfig)
	if err != nil {
		log.Printf("Error initializing RabbitMQ broker: %v", err)
		return nil, nil, err
	}

	// Crear dispatcher de eventos
	dispatcher := NewEventDispatcher(broker)

	return dispatcher, broker, nil
}

// InitConsumers inicializa los consumidores de eventos
func InitConsumers(broker EventBroker, handlers map[string]map[string]EventHandler) error {
	// Para cada topic
	for topic, eventHandlers := range handlers {
		// Crear un consumidor para el topic
		consumer := NewConsumer(broker, topic)

		// Registrar los handlers para cada tipo de evento
		for eventType, handler := range eventHandlers {
			consumer.RegisterHandler(eventType, handler)
		}

		// Iniciar el consumidor
		if err := consumer.Start(); err != nil {
			return err
		}
	}

	return nil
}
