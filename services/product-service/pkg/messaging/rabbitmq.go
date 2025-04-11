// pkg/messaging/rabbitmq.go
package messaging

import (
	"github.com/streadway/amqp"

	"github.com/ntdt/product-service/config"
)

type RabbitMQClient interface {
	Publish(exchange, routingKey string, message []byte) error
	Subscribe(exchange, routingKey, queueName string, handler func([]byte) error) error
	Close() error
}

type rabbitMQClient struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient(cfg config.RabbitMQConfig) (RabbitMQClient, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		cfg.Exchange, // exchange name
		"topic",      // exchange type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &rabbitMQClient{
		conn:     conn,
		channel:  ch,
		exchange: cfg.Exchange,
	}, nil
}

func (r *rabbitMQClient) Publish(exchange, routingKey string, message []byte) error {
	if exchange == "" {
		exchange = r.exchange
	}

	return r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)
}

func (r *rabbitMQClient) Subscribe(exchange, routingKey, queueName string, handler func([]byte) error) error {
	if exchange == "" {
		exchange = r.exchange
	}

	// Declare a queue
	q, err := r.channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange
	err = r.channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	// Consume messages
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				// Failed to process message, nack
				d.Nack(false, true)
			} else {
				// Successfully processed message, ack
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (r *rabbitMQClient) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}
