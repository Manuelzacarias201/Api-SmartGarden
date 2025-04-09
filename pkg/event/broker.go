package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// EventBroker es la interfaz para el sistema de mensajería
type EventBroker interface {
	// Publish publica un evento en un topic específico
	Publish(ctx context.Context, topic string, event Event) error

	// Subscribe suscribe a un consumidor a un topic específico
	Subscribe(topic string, handler EventHandler) error

	// Close cierra la conexión con el broker
	Close() error
}

// RabbitMQBroker implementa EventBroker usando RabbitMQ
type RabbitMQBroker struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queueName    string
}

// RabbitMQConfig contiene la configuración para conectar a RabbitMQ
type RabbitMQConfig struct {
	URL          string
	ExchangeName string
	QueueName    string
}

// NewRabbitMQBroker crea una nueva instancia de RabbitMQBroker
func NewRabbitMQBroker(config RabbitMQConfig) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declarar exchange de tipo topic
	err = ch.ExchangeDeclare(
		config.ExchangeName, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	return &RabbitMQBroker{
		conn:         conn,
		channel:      ch,
		exchangeName: config.ExchangeName,
		queueName:    config.QueueName,
	}, nil
}

// Publish publica un evento en un topic específico
func (b *RabbitMQBroker) Publish(ctx context.Context, topic string, event Event) error {
	// Convertir el evento a JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publicar el mensaje
	err = b.channel.Publish(
		b.exchangeName, // exchange
		topic,          // routing key (topic)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Mensaje persistente
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"event_type": event.Type,
				"event_id":   event.ID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Event published: %s - %s", topic, event.ID)
	return nil
}

// Subscribe suscribe a un consumidor a un topic específico
func (b *RabbitMQBroker) Subscribe(topic string, handler EventHandler) error {
	// Crear una cola para el suscriptor
	queueName := fmt.Sprintf("%s-%s", b.queueName, topic)

	q, err := b.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Enlazar la cola al exchange con el topic específico
	err = b.channel.QueueBind(
		q.Name,         // queue name
		topic,          // routing key
		b.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}

	// Consumir mensajes
	msgs, err := b.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (false para confirmar manualmente)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Procesar mensajes en una goroutine
	go func() {
		for d := range msgs {
			var event Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				d.Nack(false, true) // Rechazar y volver a la cola
				continue
			}

			// Procesar el evento con el handler
			err = handler.Handle(context.Background(), event)
			if err != nil {
				log.Printf("Error handling event: %v", err)
				d.Nack(false, true) // Rechazar y volver a la cola
				continue
			}

			// Confirmar el procesamiento exitoso
			d.Ack(false)
			log.Printf("Event processed: %s - %s", topic, event.ID)
		}
	}()

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Close cierra la conexión con RabbitMQ
func (b *RabbitMQBroker) Close() error {
	if err := b.channel.Close(); err != nil {
		return err
	}
	return b.conn.Close()
}
