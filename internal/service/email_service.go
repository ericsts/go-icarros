package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	Host string
	Port string
	From string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host: getEnvEmail("SMTP_HOST", "localhost"),
		Port: getEnvEmail("SMTP_PORT", "1025"),
		From: getEnvEmail("SMTP_FROM", "noreply@icarros.com"),
	}
}

func (e *EmailService) Send(to, subject, body string) error {
	addr := e.Host + ":" + e.Port
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		e.From, to, subject, body,
	))
	return smtp.SendMail(addr, nil, e.From, []string{to}, msg)
}

func getEnvEmail(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
