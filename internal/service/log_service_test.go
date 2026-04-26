package service

import (
	"errors"
	"testing"

	"go-icarros/internal/models"
)

type mockLogRepo struct {
	createErr     error
	findAllResult []models.EventLog
	findAllErr    error
	lastEntry     *models.EventLog
}

func (m *mockLogRepo) Create(entry *models.EventLog) error {
	m.lastEntry = entry
	return m.createErr
}

func (m *mockLogRepo) FindAll(_, _ string, _ int) ([]models.EventLog, error) {
	return m.findAllResult, m.findAllErr
}

func TestLogService_Info(t *testing.T) {
	repo := &mockLogRepo{}
	svc := &LogService{Repo: repo}

	svc.Info("car.created", "carro cadastrado", map[string]any{"car_id": 1})

	if repo.lastEntry == nil {
		t.Fatal("esperado registro salvo no repo")
	}
	if repo.lastEntry.Level != "info" {
		t.Errorf("esperado level=info, obtido %s", repo.lastEntry.Level)
	}
	if repo.lastEntry.Event != "car.created" {
		t.Errorf("esperado event=car.created, obtido %s", repo.lastEntry.Event)
	}
	if repo.lastEntry.Message != "carro cadastrado" {
		t.Errorf("esperado message='carro cadastrado', obtido %s", repo.lastEntry.Message)
	}
}

func TestLogService_Warn(t *testing.T) {
	repo := &mockLogRepo{}
	svc := &LogService{Repo: repo}

	svc.Warn("auction.create_failed", "falha ao criar leilão", nil)

	if repo.lastEntry == nil {
		t.Fatal("esperado registro salvo no repo")
	}
	if repo.lastEntry.Level != "warn" {
		t.Errorf("esperado level=warn, obtido %s", repo.lastEntry.Level)
	}
}

func TestLogService_Error(t *testing.T) {
	repo := &mockLogRepo{}
	svc := &LogService{Repo: repo}

	svc.Error("notification.email", "falha ao enviar", map[string]any{"error": "timeout"})

	if repo.lastEntry == nil {
		t.Fatal("esperado registro salvo no repo")
	}
	if repo.lastEntry.Level != "error" {
		t.Errorf("esperado level=error, obtido %s", repo.lastEntry.Level)
	}
}

func TestLogService_GetAll(t *testing.T) {
	logs := []models.EventLog{{ID: 1, Level: "info"}, {ID: 2, Level: "error"}}
	svc := &LogService{Repo: &mockLogRepo{findAllResult: logs}}

	result, err := svc.GetAll("", "", 100)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("esperado 2 logs, obtido %d", len(result))
	}
}

func TestLogService_GetAll_PropagaErro(t *testing.T) {
	svc := &LogService{Repo: &mockLogRepo{findAllErr: errors.New("db error")}}

	_, err := svc.GetAll("", "", 100)

	if err == nil {
		t.Fatal("deveria retornar erro")
	}
}
