package email

import (
	"encoding/json"
	"mailer-rabbitmq/internal/models"
	"mailer-rabbitmq/internal/service"
	"net/http"
)

type EmailRequest struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To       []models.EmailRecipient `json:"to"`
	Subject  string                  `json:"subject"`
	Text     string                  `json:"text"`
	Category string                  `json:"category"`
}

type Handler struct {
	service *service.EmailService
}

func NewHandler(service *service.EmailService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var emailReq EmailRequest

	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := models.Email{
		From:     emailReq.From,
		To:       emailReq.To,
		Subject:  emailReq.Subject,
		Text:     emailReq.Text,
		Category: emailReq.Category,
	}

	if err := h.service.SendEmail(email); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
