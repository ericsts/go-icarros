package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-icarros/internal/service"

	"github.com/gin-gonic/gin"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware_SemHeader(t *testing.T) {
	r := newTestRouter()
	r.GET("/protegido", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtido %d", w.Code)
	}
}

func TestAuthMiddleware_TokenInvalido(t *testing.T) {
	r := newTestRouter()
	r.GET("/protegido", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Bearer tokeninvalido")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtido %d", w.Code)
	}
}

func TestAuthMiddleware_TokenValido(t *testing.T) {
	token, _ := service.GenerateToken(1, "user")

	r := newTestRouter()
	r.GET("/protegido", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestAuthMiddleware_DefineUserIDERole(t *testing.T) {
	token, _ := service.GenerateToken(42, "admin")

	var gotUserID, gotRole any

	r := newTestRouter()
	r.GET("/protegido", AuthMiddleware(), func(c *gin.Context) {
		gotUserID, _ = c.Get("user_id")
		gotRole, _ = c.Get("role")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if gotUserID != 42 {
		t.Errorf("esperado user_id=42, obtido %v", gotUserID)
	}
	if gotRole != "admin" {
		t.Errorf("esperado role=admin, obtido %v", gotRole)
	}
}

func TestAuthMiddleware_SemPrefixoBearer(t *testing.T) {
	token, _ := service.GenerateToken(1, "user")

	r := newTestRouter()
	r.GET("/protegido", AuthMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// token sem "Bearer " → deve rejeitar
	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401 sem prefixo Bearer, obtido %d", w.Code)
	}
}
