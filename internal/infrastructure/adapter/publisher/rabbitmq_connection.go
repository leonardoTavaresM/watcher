package publisher

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMQConnection(uri string) (*RabbitMQConnection, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return &RabbitMQConnection{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *RabbitMQConnection) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}

	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

func (c *RabbitMQConnection) GetChannel() *amqp091.Channel {
	return c.channel
}
