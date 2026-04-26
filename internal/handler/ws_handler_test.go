package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-icarros/internal/service"
	"go-icarros/internal/ws"
)

func TestWSHandler_ServeAuction_IDInvalido(t *testing.T) {
	h := &WSHandler{Hub: ws.NewHub()}
	r := newTestRouter()
	r.GET("/ws/auctions/:id", h.ServeAuction)

	req := httptest.NewRequest(http.MethodGet, "/ws/auctions/abc?token=qualquer", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtido %d", w.Code)
	}
}

func TestWSHandler_ServeAuction_TokenAusente(t *testing.T) {
	h := &WSHandler{Hub: ws.NewHub()}
	r := newTestRouter()
	r.GET("/ws/auctions/:id", h.ServeAuction)

	req := httptest.NewRequest(http.MethodGet, "/ws/auctions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtido %d", w.Code)
	}
}

func TestWSHandler_ServeAuction_TokenInvalido(t *testing.T) {
	h := &WSHandler{Hub: ws.NewHub()}
	r := newTestRouter()
	r.GET("/ws/auctions/:id", h.ServeAuction)

	req := httptest.NewRequest(http.MethodGet, "/ws/auctions/1?token=isso.nao.e.um.jwt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtido %d", w.Code)
	}
}

func TestWSHandler_ServeAuction_TokenExpirado(t *testing.T) {
	// token assinado com chave errada simula token adulterado
	h := &WSHandler{Hub: ws.NewHub()}
	r := newTestRouter()
	r.GET("/ws/auctions/:id", h.ServeAuction)

	// JWT válido estruturalmente mas assinado com chave diferente
	tokenFalso := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImV4cCI6MTYwMDAwMDAwMH0.assinatura_invalida"
	req := httptest.NewRequest(http.MethodGet, "/ws/auctions/1?token="+tokenFalso, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperado 401, obtido %d", w.Code)
	}
}

func TestWSHandler_ServeAuction_TokenValido_UpgradeNecessario(t *testing.T) {
	// Com token válido, a lógica de auth passa e o handler tenta o upgrade WebSocket.
	// httptest.ResponseRecorder não suporta hijack, então o upgrader falha e retorna
	// um erro HTTP — mas o importante é que NÃO retorna 401/403 (auth funcionou).
	token, err := service.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("erro ao gerar token: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	h := &WSHandler{Hub: hub}
	r := newTestRouter()
	r.GET("/ws/auctions/:id", h.ServeAuction)

	req := httptest.NewRequest(http.MethodGet, "/ws/auctions/1?token="+token, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden {
		t.Errorf("token válido não deve retornar %d", w.Code)
	}
}
