package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-icarros/internal/models"

	"github.com/gin-gonic/gin"
)

type mockCarSvc struct {
	createErr       error
	getAllCars      []models.Car
	getAllErr       error
	getByIDCar      *models.Car
	getByIDErr      error
	getByUserIDCars []models.Car
	getByUserIDErr  error
	updateErr       error
	deleteErr       error
}

func (m *mockCarSvc) Create(_ *models.Car) error         { return m.createErr }
func (m *mockCarSvc) GetAll() ([]models.Car, error)      { return m.getAllCars, m.getAllErr }
func (m *mockCarSvc) GetByID(_ int) (*models.Car, error) { return m.getByIDCar, m.getByIDErr }
func (m *mockCarSvc) GetByUserID(_ int) ([]models.Car, error) {
	return m.getByUserIDCars, m.getByUserIDErr
}
func (m *mockCarSvc) Update(_ *models.Car) error { return m.updateErr }
func (m *mockCarSvc) Delete(_ int) error         { return m.deleteErr }

// injetaUserID simula o middleware de autenticação nos testes
func injetaUserID(userID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func TestCarHandler_Create_Sucesso(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{}}
	r := newTestRouter()
	r.POST("/cars", injetaUserID(1), h.Create)

	body := models.Car{Marca: "VW", Modelo: "Gol", Ano: 2020, Valor: 45000}
	req := httptest.NewRequest(http.MethodPost, "/cars", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtido %d: %s", w.Code, w.Body.String())
	}
}

func TestCarHandler_Create_DefineUserID(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{}}
	r := newTestRouter()
	r.POST("/cars", injetaUserID(7), h.Create)

	req := httptest.NewRequest(http.MethodPost, "/cars", jsonBody(models.Car{Marca: "Fiat"}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var car models.Car
	json.Unmarshal(w.Body.Bytes(), &car)
	if car.UserID != 7 {
		t.Errorf("esperado UserID=7, obtido %d", car.UserID)
	}
}

func TestCarHandler_Create_ErroService(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{createErr: errors.New("db error")}}
	r := newTestRouter()
	r.POST("/cars", injetaUserID(1), h.Create)

	req := httptest.NewRequest(http.MethodPost, "/cars", jsonBody(models.Car{Marca: "VW"}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtido %d", w.Code)
	}
}

func TestCarHandler_List(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{getAllCars: []models.Car{{ID: 1}, {ID: 2}}}}
	r := newTestRouter()
	r.GET("/cars", h.List)

	req := httptest.NewRequest(http.MethodGet, "/cars", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
	var cars []models.Car
	json.Unmarshal(w.Body.Bytes(), &cars)
	if len(cars) != 2 {
		t.Errorf("esperado 2 carros, obtido %d", len(cars))
	}
}

func TestCarHandler_GetByID_Sucesso(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{getByIDCar: &models.Car{ID: 1, Marca: "Fiat"}}}
	r := newTestRouter()
	r.GET("/cars/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/cars/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestCarHandler_GetByID_NaoEncontrado(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{getByIDErr: errors.New("not found")}}
	r := newTestRouter()
	r.GET("/cars/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/cars/99", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtido %d", w.Code)
	}
}

func TestCarHandler_GetMyCars(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{getByUserIDCars: []models.Car{{ID: 1}, {ID: 2}}}}
	r := newTestRouter()
	r.GET("/cars/my", injetaUserID(1), h.GetMyCars)

	req := httptest.NewRequest(http.MethodGet, "/cars/my", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
	var cars []models.Car
	json.Unmarshal(w.Body.Bytes(), &cars)
	if len(cars) != 2 {
		t.Errorf("esperado 2 carros, obtido %d", len(cars))
	}
}

func TestCarHandler_Update(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{}}
	r := newTestRouter()
	r.PUT("/cars/:id", h.Update)

	body := models.Car{Marca: "Fiat", Modelo: "Uno", Ano: 2021, Valor: 30000}
	req := httptest.NewRequest(http.MethodPut, "/cars/1", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestCarHandler_Delete(t *testing.T) {
	h := &CarHandler{Service: &mockCarSvc{}}
	r := newTestRouter()
	r.DELETE("/cars/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/cars/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperado 204, obtido %d", w.Code)
	}
}
