package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/leonardoTavaresM/watcher/internal/domain/entity"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn       *RabbitMQConnection
	routingKey string
	exchange   string
}

// NewRabbitMQPublisher NÃO declara nem faz bind de fila nenhuma: um producer
// puro só precisa que o exchange exista. Declarar/consumir filas é
// responsabilidade de quem consome (ver serviço collector) — do contrário o
// watcher fica "dono" de uma fila que ele mesmo nunca lê, e ela cresce sem
// limite (foi exatamente esse bug que existia aqui antes desse ajuste).
func NewRabbitMQPublisher(uri, exchange, routingKey string) (*RabbitMQPublisher, error) {
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

	return &RabbitMQPublisher{
		conn:       conn,
		routingKey: routingKey,
		exchange:   exchange,
	}, nil
}

func (p *RabbitMQPublisher) Publish(event entity.FileEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.conn.GetChannel().Publish(
		p.exchange,   // exchange
		p.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
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
