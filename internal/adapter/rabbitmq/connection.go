package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// Connect to RabbitMQ
func NewConnection(uri string) (*Connection, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return &Connection{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close connection

func (c *Connection) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}

	if c.conn != nil {
		c.channel.Close()
	}
	return nil
}

func (c *Connection) GetChannel() *amqp091.Channel {
	return c.channel
}
