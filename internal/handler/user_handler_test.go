package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-icarros/internal/models"

	"github.com/gin-gonic/gin"
)

type mockUserSvc struct {
	registerErr error
	loginUser   *models.User
	loginErr    error
	getAllUsers []models.User
	getAllErr   error
	getByIDUser *models.User
	getByIDErr  error
	updateErr   error
	deleteErr   error
}

func (m *mockUserSvc) Register(_ *models.User) error           { return m.registerErr }
func (m *mockUserSvc) Login(_, _ string) (*models.User, error) { return m.loginUser, m.loginErr }
func (m *mockUserSvc) GetAll() ([]models.User, error)          { return m.getAllUsers, m.getAllErr }
func (m *mockUserSvc) GetByID(_ int) (*models.User, error)     { return m.getByIDUser, m.getByIDErr }
func (m *mockUserSvc) Update(_ *models.User) error             { return m.updateErr }
func (m *mockUserSvc) Delete(_ int) error                      { return m.deleteErr }

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func jsonBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func TestUserHandler_Login_Sucesso(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{loginUser: &models.User{ID: 1, Role: "user"}}}
	r := newTestRouter()
	r.POST("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(map[string]string{"email": "a@a.com", "password": "123"}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Error("esperado token na resposta")
	}
}

func TestUserHandler_Login_CredenciaisInvalidas(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{loginErr: errors.New("invalid")}}
	r := newTestRouter()
	r.POST("/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(map[string]string{"email": "a@a.com", "password": "errada"}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtido %d", w.Code)
	}
}

func TestUserHandler_Create_Sucesso(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{}}
	r := newTestRouter()
	r.POST("/users", h.Create)

	body := models.User{Name: "Eric", Email: "eric@test.com", Password: "123", Role: "user"}
	req := httptest.NewRequest(http.MethodPost, "/users", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtido %d: %s", w.Code, w.Body.String())
	}
}

func TestUserHandler_Create_SenhaOcultaNaResposta(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{}}
	r := newTestRouter()
	r.POST("/users", h.Create)

	body := models.User{Name: "Eric", Email: "eric@test.com", Password: "secreta", Role: "user"}
	req := httptest.NewRequest(http.MethodPost, "/users", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if _, ok := resp["password"]; ok {
		t.Error("senha não deve aparecer na resposta")
	}
}

func TestUserHandler_Create_ErroService(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{registerErr: errors.New("db error")}}
	r := newTestRouter()
	r.POST("/users", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/users", jsonBody(models.User{Name: "x", Password: "x", Role: "user"}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtido %d", w.Code)
	}
}

func TestUserHandler_List(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{getAllUsers: []models.User{{ID: 1}, {ID: 2}}}}
	r := newTestRouter()
	r.GET("/users", h.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
	var users []models.User
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 2 {
		t.Errorf("esperado 2 usuários, obtido %d", len(users))
	}
}

func TestUserHandler_GetByID_Sucesso(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{getByIDUser: &models.User{ID: 1, Name: "Eric"}}}
	r := newTestRouter()
	r.GET("/users/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestUserHandler_GetByID_NaoEncontrado(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{getByIDErr: errors.New("not found")}}
	r := newTestRouter()
	r.GET("/users/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtido %d", w.Code)
	}
}

func TestUserHandler_GetByID_IDInvalido(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{}}
	r := newTestRouter()
	r.GET("/users/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtido %d", w.Code)
	}
}

func TestUserHandler_Update(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{}}
	r := newTestRouter()
	r.PUT("/users/:id", h.Update)

	body := models.User{Name: "Atualizado", Email: "u@u.com", Role: "user"}
	req := httptest.NewRequest(http.MethodPut, "/users/1", jsonBody(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestUserHandler_Delete(t *testing.T) {
	h := &UserHandler{Service: &mockUserSvc{}}
	r := newTestRouter()
	r.DELETE("/users/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("esperado 204, obtido %d", w.Code)
	}
}
