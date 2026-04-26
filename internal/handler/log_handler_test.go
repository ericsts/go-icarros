package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-icarros/internal/models"
)

type mockLogSvc struct {
	getAllResult []models.EventLog
	getAllErr    error
}

func (m *mockLogSvc) GetAll(_, _ string, _ int) ([]models.EventLog, error) {
	return m.getAllResult, m.getAllErr
}

func TestLogHandler_List(t *testing.T) {
	logs := []models.EventLog{
		{ID: 1, Level: "info", Event: "car.created", Message: "carro cadastrado"},
		{ID: 2, Level: "error", Event: "notification.email", Message: "falha ao enviar"},
	}
	h := &LogHandler{Service: &mockLogSvc{getAllResult: logs}}
	r := newTestRouter()
	r.GET("/logs", h.List)

	req := httptest.NewRequest(http.MethodGet, "/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
	var result []models.EventLog
	json.Unmarshal(w.Body.Bytes(), &result)
	if len(result) != 2 {
		t.Errorf("esperado 2 logs, obtido %d", len(result))
	}
}

func TestLogHandler_List_ComFiltros(t *testing.T) {
	logs := []models.EventLog{{ID: 1, Level: "error", Event: "notification.email"}}
	h := &LogHandler{Service: &mockLogSvc{getAllResult: logs}}
	r := newTestRouter()
	r.GET("/logs", h.List)

	req := httptest.NewRequest(http.MethodGet, "/logs?level=error&event=notification&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestLogHandler_List_Erro(t *testing.T) {
	h := &LogHandler{Service: &mockLogSvc{getAllErr: errors.New("db error")}}
	r := newTestRouter()
	r.GET("/logs", h.List)

	req := httptest.NewRequest(http.MethodGet, "/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtido %d", w.Code)
	}
}
