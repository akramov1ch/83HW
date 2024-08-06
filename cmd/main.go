package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	email "mailer-rabbitmq/internal/handler"
	"mailer-rabbitmq/internal/mailer"
	"mailer-rabbitmq/internal/models"
	"mailer-rabbitmq/internal/producer"
	"mailer-rabbitmq/internal/rabbitmq"
	"mailer-rabbitmq/internal/service"

	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	rabbitMQConfig := rabbitmq.Config{
		URI:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "email_exchange",
		ExchangeType: "direct",
		Queue:        "email_queue",
		RoutingKey:   "email_key",
	}

	rabbitMQPublisher, err := producer.NewPublisher(rabbitMQConfig)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %s", err)
	}
	defer rabbitMQPublisher.Close()

	service := service.NewEmailService(rabbitMQPublisher)
	handler := email.NewHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/send-email", handler.SendEmail).Methods("POST")

	http.Handle("/", r)

	log.Println("Starting server on ", os.Getenv("SERVER_PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

func sendViaAPI() error {
	url := "https://sandbox.api.mailtrap.io/api/send/2319965"

	payload := models.Email{
		From: models.EmailAddress{
			Email: "mailtrap@example.com",
			Name:  "Mailtrap Test",
		},
		To: []models.EmailRecipient{
			{Email: "firdavs792@gmail.com"},
		},
		Subject:  "You are awesome! v2",
		Text:     "Congrats for sending test email with Mailtrap!",
		Category: "Integration Test",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("EMAIL_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}

func sendViaSMTP(recipient string) error {
	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("failed to convert 'SMTP_PORT' to int: %w", err)
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	sender := os.Getenv("SMTP_SENDER")

	m := mailer.New(host, port, username, password, sender)

	data := map[string]any{
		"ID": rand.Intn(100),
	}

	err = m.Send(recipient, "welcome.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Successfully send an email")
	return nil
}
