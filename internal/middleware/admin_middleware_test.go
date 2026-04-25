package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminMiddleware_SemRole(t *testing.T) {
	r := newTestRouter()
	r.GET("/admin", AdminMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado 403, obtido %d", w.Code)
	}
}

func TestAdminMiddleware_RoleUser(t *testing.T) {
	r := newTestRouter()
	r.GET("/admin",
		func(c *gin.Context) { c.Set("role", "user"); c.Next() },
		AdminMiddleware(),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("esperado 403 para role=user, obtido %d", w.Code)
	}
}

func TestAdminMiddleware_RoleAdmin(t *testing.T) {
	r := newTestRouter()
	r.GET("/admin",
		func(c *gin.Context) { c.Set("role", "admin"); c.Next() },
		AdminMiddleware(),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200 para role=admin, obtido %d", w.Code)
	}
}
