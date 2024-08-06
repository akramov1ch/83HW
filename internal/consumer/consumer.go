package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"mailer-rabbitmq/internal/models"
	"mailer-rabbitmq/internal/producer"
	"mailer-rabbitmq/internal/rabbitmq"
	"mailer-rabbitmq/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func StartConsumer() {
    conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "email_queue", // name
        true,          // durable
        false,         // delete when unused
        false,         // exclusive
        false,         // no-wait
        nil,           // arguments
    )
    failOnError(err, "Failed to declare a queue")

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    publisher, err  := producer.NewPublisher(rabbitmq.Config{})
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}
    emailService := service.NewEmailService(publisher)

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
            
            var emailMsg models.Email
            err := json.Unmarshal(d.Body, &emailMsg)
            if err != nil {
                log.Printf("Error decoding JSON: %s", err)
                continue
            }

            err = emailService.SendEmail(emailMsg)
            if err != nil {
                log.Printf("Error sending email: %s", err)
            } else {
                log.Printf("Email sent to: %s", emailMsg.To)
            }
        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}
