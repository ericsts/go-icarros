package repository

import (
	"encoding/json"
	"testing"
	"time"

	"go-icarros/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLogRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO event_logs").
		WithArgs("info", "car.created", "carro cadastrado", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &LogRepository{DB: db}
	err = repo.Create(&models.EventLog{
		Level:    "info",
		Event:    "car.created",
		Message:  "carro cadastrado",
		Metadata: map[string]any{"car_id": 1},
	})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestLogRepository_FindAll_SemFiltros(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	meta, _ := json.Marshal(map[string]any{"car_id": 1})
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "level", "event", "message", "metadata", "created_at"}).
		AddRow(1, "info", "car.created", "carro cadastrado", meta, now).
		AddRow(2, "error", "notification.email", "falha ao enviar", meta, now)

	mock.ExpectQuery("SELECT id, level, event, message, metadata, created_at FROM event_logs").
		WithArgs(100).
		WillReturnRows(rows)

	repo := &LogRepository{DB: db}
	logs, err := repo.FindAll("", "", 100)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("esperado 2 logs, obtido %d", len(logs))
	}
	if logs[0].Level != "info" {
		t.Errorf("esperado level=info, obtido %s", logs[0].Level)
	}
}

func TestLogRepository_FindAll_FiltroLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	meta, _ := json.Marshal(map[string]any{})
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "level", "event", "message", "metadata", "created_at"}).
		AddRow(1, "error", "notification.email", "falha ao enviar", meta, now)

	mock.ExpectQuery("SELECT id, level, event, message, metadata, created_at FROM event_logs").
		WithArgs("error", 100).
		WillReturnRows(rows)

	repo := &LogRepository{DB: db}
	logs, err := repo.FindAll("error", "", 100)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("esperado 1 log, obtido %d", len(logs))
	}
	if logs[0].Level != "error" {
		t.Errorf("esperado level=error, obtido %s", logs[0].Level)
	}
}

func TestLogRepository_FindAll_FiltroEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	meta, _ := json.Marshal(map[string]any{})
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "level", "event", "message", "metadata", "created_at"}).
		AddRow(1, "info", "car.created", "carro cadastrado", meta, now)

	mock.ExpectQuery("SELECT id, level, event, message, metadata, created_at FROM event_logs").
		WithArgs("%car.created%", 100).
		WillReturnRows(rows)

	repo := &LogRepository{DB: db}
	logs, err := repo.FindAll("", "car.created", 100)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("esperado 1 log, obtido %d", len(logs))
	}
	if logs[0].Event != "car.created" {
		t.Errorf("esperado event=car.created, obtido %s", logs[0].Event)
	}
}

func TestLogRepository_FindAll_TodosFiltros(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	meta, _ := json.Marshal(map[string]any{})
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "level", "event", "message", "metadata", "created_at"}).
		AddRow(1, "info", "car.created", "carro cadastrado", meta, now)

	mock.ExpectQuery("SELECT id, level, event, message, metadata, created_at FROM event_logs").
		WithArgs("info", "%car%", 10).
		WillReturnRows(rows)

	repo := &LogRepository{DB: db}
	logs, err := repo.FindAll("info", "car", 10)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("esperado 1 log, obtido %d", len(logs))
	}
}
