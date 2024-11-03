package queue

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func CreateQueueClient() (*RabbitMQClient, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, ErrEnvVarNotSet
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, ErrConnFailed
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, ErrChanFailed
	}

	return &RabbitMQClient{conn: conn, ch: ch}, nil
}

func (r *RabbitMQClient) Reconnect() error {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return ErrEnvVarNotSet
	}

	if r.conn == nil || r.conn.IsClosed() {
		log.Println("Connection to RabbitMQ is closed. Reconnecting...")
		conn, err := amqp.Dial(url)
		if err != nil {
			return ErrConnFailed
		}
		r.conn = conn
	}

	if r.ch == nil {
		log.Println("Channel to RabbitMQ is closed. Reconnecting...")
		ch, err := r.conn.Channel()
		if err != nil {
			return ErrChanFailed
		}
		r.ch = ch
	}

	return nil
}

func (r *RabbitMQClient) PublishMessage(msgContext []byte, queueName string) error {
	if err := r.Reconnect(); err != nil {
		return err
	}

	err := r.ch.Publish(
		"",        // default exchange
		queueName, // routing key (queue name)
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgContext,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	return err
}

func (r *RabbitMQClient) Close() {
	r.ch.Close()
	r.conn.Close()
}
