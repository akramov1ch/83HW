package service

import (
	"log"
	"mailer-rabbitmq/internal/models"
	"mailer-rabbitmq/internal/producer"
)

type EmailService struct {
	Publisher *producer.Publisher
}

func NewEmailService(publisher *producer.Publisher) *EmailService {
	return &EmailService{
		Publisher: publisher,
	}
}

func (s *EmailService) SendEmail(email models.Email) error {
	log.Printf("Preparing to send email to %v", email.To)

	if err := s.Publisher.Publish(email); err != nil {
		return err
	}

	return nil
}
