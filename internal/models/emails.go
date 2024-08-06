package models

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type EmailRecipient struct {
	Email string `json:"email"`
}

type Email struct {
	From     EmailAddress     `json:"from"`
	To       []EmailRecipient `json:"to"`
	Subject  string           `json:"subject"`
	Text     string           `json:"text"`
	Category string           `json:"category"`
}
