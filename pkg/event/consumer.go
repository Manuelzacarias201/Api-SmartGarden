package event

import (
	"context"
	"log"
)

// Consumer es un consumidor de eventos que procesa los eventos de un topic específico
type Consumer struct {
	broker   EventBroker
	topic    string
	handlers map[string]EventHandler
}

// NewConsumer crea un nuevo consumidor
func NewConsumer(broker EventBroker, topic string) *Consumer {
	return &Consumer{
		broker:   broker,
		topic:    topic,
		handlers: make(map[string]EventHandler),
	}
}

// RegisterHandler registra un manejador para un tipo de evento específico
func (c *Consumer) RegisterHandler(eventType string, handler EventHandler) {
	c.handlers[eventType] = handler
}

// Start inicia el consumidor
func (c *Consumer) Start() error {
	return c.broker.Subscribe(c.topic, c)
}

// Handle implementa la interfaz EventHandler para procesar eventos
func (c *Consumer) Handle(ctx context.Context, event Event) error {
	handler, exists := c.handlers[event.Type]
	if !exists {
		log.Printf("No handler registered for event type: %s", event.Type)
		return nil
	}

	return handler.Handle(ctx, event)
}
