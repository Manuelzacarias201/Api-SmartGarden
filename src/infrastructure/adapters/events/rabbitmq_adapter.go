package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ApiSmart/src/core/application"
	"ApiSmart/src/core/domain/events"
	"ApiSmart/src/core/tipo_de_datos"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// RabbitMQAdapter implementa application.EventBroker usando RabbitMQ
type RabbitMQAdapter struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queueName    string
}

// NewRabbitMQAdapter crea una nueva instancia de RabbitMQAdapter
func NewRabbitMQAdapter(config tipo_de_datos.RabbitMQConfig) (*RabbitMQAdapter, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %w", err)
	}

	// Abrir un canal
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error abriendo canal: %w", err)
	}

	// Declarar exchange
	err = ch.ExchangeDeclare(
		config.ExchangeName, // nombre
		"topic",             // tipo
		true,                // durable
		false,               // auto-delete
		false,               // internal
		false,               // no-wait
		nil,                 // argumentos
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("error declarando exchange: %w", err)
	}

	return &RabbitMQAdapter{
		conn:         conn,
		channel:      ch,
		exchangeName: config.ExchangeName,
		queueName:    config.QueueName,
	}, nil
}

// Close cierra la conexión con RabbitMQ
func (a *RabbitMQAdapter) Close() error {
	if err := a.channel.Close(); err != nil {
		return err
	}
	return a.conn.Close()
}

// Publish publica un evento en RabbitMQ
func (a *RabbitMQAdapter) Publish(ctx context.Context, topic string, event events.Event) error {
	// Convertir el evento a JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando evento: %w", err)
	}

	// Publicar mensaje
	err = a.channel.Publish(
		a.exchangeName, // exchange
		topic,          // routing key (topic)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"event_type": event.Type,
				"event_id":   event.ID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error publicando mensaje: %w", err)
	}

	log.Printf("Evento publicado: %s - %s", topic, event.ID)
	return nil
}

// Subscribe suscribe a un topic de RabbitMQ
func (a *RabbitMQAdapter) Subscribe(topic string, handler application.EventHandler) error {
	// Crear una cola específica para este consumidor
	queueName := fmt.Sprintf("%s-%s", a.queueName, topic)

	q, err := a.channel.QueueDeclare(
		queueName, // nombre
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("error declarando cola: %w", err)
	}

	// Enlazar la cola al exchange con el routing key (topic)
	err = a.channel.QueueBind(
		q.Name,         // nombre de cola
		topic,          // routing key
		a.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error enlazando cola: %w", err)
	}

	// Consumir mensajes
	msgs, err := a.channel.Consume(
		q.Name, // cola
		"",     // consumer
		false,  // auto-ack (false para confirmar manualmente)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("error registrando consumidor: %w", err)
	}

	// Procesar mensajes en una goroutine
	go func() {
		for d := range msgs {
			var event events.Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Error deserializando evento: %v", err)
				d.Nack(false, true) // Rechazar y volver a la cola
				continue
			}

			// Procesar el evento
			err = handler.Handle(context.Background(), event)
			if err != nil {
				log.Printf("Error procesando evento: %v", err)
				d.Nack(false, true) // Rechazar y volver a la cola
				continue
			}

			d.Ack(false) // Confirmar procesamiento
			log.Printf("Evento procesado: %s - %s", topic, event.ID)
		}
	}()

	log.Printf("Suscrito al topic: %s", topic)
	return nil
}

// EventDispatcherAdapter implementa application.EventDispatcher
type EventDispatcherAdapter struct {
	broker application.EventBroker
}

// NewEventDispatcherAdapter crea una nueva instancia de EventDispatcherAdapter
func NewEventDispatcherAdapter(broker application.EventBroker) *EventDispatcherAdapter {
	return &EventDispatcherAdapter{
		broker: broker,
	}
}

// Dispatch envía un evento al broker
func (d *EventDispatcherAdapter) Dispatch(ctx context.Context, eventType string, topic string, data map[string]interface{}) error {
	event := events.Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	return d.broker.Publish(ctx, topic, event)
}
