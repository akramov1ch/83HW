package producer

import (
	"encoding/json"
	"log"
	"mailer-rabbitmq/internal/models"
	"mailer-rabbitmq/internal/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  rabbitmq.Config
}

func NewPublisher(config rabbitmq.Config) (*Publisher, error) {
	conn, err := amqp.Dial(config.URI)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		config.Exchange,     // name of the exchange
		config.ExchangeType, // type of the exchange
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

func (publisher *Publisher) Publish(email models.Email) error {
	body, err := json.Marshal(email)
	if err != nil {
		return err
	}

	err = publisher.channel.Publish(
		publisher.config.Exchange,   // exchange
		publisher.config.RoutingKey, // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Sent email to RabbitMQ: %s", body)
	return nil
}

func (publisher *Publisher) Close() {
	if err := publisher.channel.Close(); err != nil {
		log.Fatalf("Failed to close channel: %s", err)
	}
	if err := publisher.conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %s", err)
	}
}
