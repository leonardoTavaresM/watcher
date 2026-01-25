package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain/entity"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn     *RabbitMQConnection
	queue    string
	exchange string
}

func NewRabbitMQPublisher(uri, exchange, queue string) (*RabbitMQPublisher, error) {
	conn, err := NewRabbitMQConnection(uri)
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = conn.GetChannel().ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err = conn.GetChannel().QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = conn.GetChannel().QueueBind(
		queue,    // queue name
		queue,    // routing key
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQPublisher{
		conn:     conn,
		queue:    queue,
		exchange: exchange,
	}, nil
}

func (p *RabbitMQPublisher) Publish(event entity.FileEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.conn.GetChannel().Publish(
		p.exchange, // exchange
		p.queue,    // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *RabbitMQPublisher) Close() error {
	return p.conn.Close()
}
