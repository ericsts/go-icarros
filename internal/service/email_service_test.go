package service

import (
	"os"
	"testing"
)

func TestNewEmailService_ValoresPadrao(t *testing.T) {
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_FROM")

	svc := NewEmailService()

	if svc.Host != "localhost" {
		t.Errorf("esperado Host=localhost, obtido %s", svc.Host)
	}
	if svc.Port != "1025" {
		t.Errorf("esperado Port=1025, obtido %s", svc.Port)
	}
	if svc.From != "noreply@icarros.com" {
		t.Errorf("esperado From=noreply@icarros.com, obtido %s", svc.From)
	}
}

func TestNewEmailService_LeDasEnvVars(t *testing.T) {
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_FROM", "contato@example.com")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_FROM")
	}()

	svc := NewEmailService()

	if svc.Host != "smtp.example.com" {
		t.Errorf("esperado Host=smtp.example.com, obtido %s", svc.Host)
	}
	if svc.Port != "587" {
		t.Errorf("esperado Port=587, obtido %s", svc.Port)
	}
	if svc.From != "contato@example.com" {
		t.Errorf("esperado From=contato@example.com, obtido %s", svc.From)
	}
}

func TestNewEmailService_NaoEhNil(t *testing.T) {
	svc := NewEmailService()
	if svc == nil {
		t.Fatal("NewEmailService não deve retornar nil")
	}
}
