package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-icarros/internal/models"
)

// mockAuctionSvc implementa AuctionSvc com campos configuráveis.
// Também é usado em car_handler_test.go (mesmo pacote).
// mockAuctionSvc implementa AuctionSvc com campos configuráveis.
// Também é usado em car_handler_test.go (mesmo pacote).
type mockAuctionSvc struct {
	createForCarResult *models.Auction
	createForCarErr    error
	hasOpenAuction     bool
	hasOpenAuctionErr  error
	getAllResult       []models.Auction
	getAllErr          error
	getByIDResult      *models.Auction
	getByIDErr         error
	placeBidResult     *models.Bid
	placeBidErr        error
	getBidsResult      []models.Bid
	getBidsErr         error
}

func (m *mockAuctionSvc) CreateForCar(_ int, _ time.Time, _ float64) (*models.Auction, error) {
	return m.createForCarResult, m.createForCarErr
}
func (m *mockAuctionSvc) HasOpenAuction(_ int) (bool, error) {
	return m.hasOpenAuction, m.hasOpenAuctionErr
}
func (m *mockAuctionSvc) GetAll() ([]models.Auction, error) {
	return m.getAllResult, m.getAllErr
}
func (m *mockAuctionSvc) GetByID(_ int) (*models.Auction, error) {
	return m.getByIDResult, m.getByIDErr
}
func (m *mockAuctionSvc) PlaceBid(_ int, _ int, _ float64) (*models.Bid, error) {
	return m.placeBidResult, m.placeBidErr
}
func (m *mockAuctionSvc) GetBids(_ int) ([]models.Bid, error) {
	return m.getBidsResult, m.getBidsErr
}

func TestAuctionHandler_List(t *testing.T) {
	auctions := []models.Auction{{ID: 1, Status: "open"}, {ID: 2, Status: "closed"}}
	h := &AuctionHandler{Service: &mockAuctionSvc{getAllResult: auctions}}
	r := newTestRouter()
	r.GET("/auctions", h.List)

	req := httptest.NewRequest(http.MethodGet, "/auctions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestAuctionHandler_List_Erro(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{getAllErr: errors.New("db error")}}
	r := newTestRouter()
	r.GET("/auctions", h.List)

	req := httptest.NewRequest(http.MethodGet, "/auctions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtido %d", w.Code)
	}
}

func TestAuctionHandler_GetByID_Sucesso(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{
		getByIDResult: &models.Auction{ID: 1, Status: "open"},
	}}
	r := newTestRouter()
	r.GET("/auctions/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/auctions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestAuctionHandler_GetByID_IDInvalido(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{}}
	r := newTestRouter()
	r.GET("/auctions/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/auctions/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtido %d", w.Code)
	}
}

func TestAuctionHandler_GetByID_NaoEncontrado(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{getByIDErr: errors.New("not found")}}
	r := newTestRouter()
	r.GET("/auctions/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/auctions/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("esperado 404, obtido %d", w.Code)
	}
}

func TestAuctionHandler_PlaceBid_Sucesso(t *testing.T) {
	bid := &models.Bid{ID: 1, AuctionID: 1, UserID: 2, Amount: 6000}
	h := &AuctionHandler{Service: &mockAuctionSvc{placeBidResult: bid}}
	r := newTestRouter()
	r.POST("/auctions/:id/bids", injetaUserID(2), h.PlaceBid)

	req := httptest.NewRequest(http.MethodPost, "/auctions/1/bids", jsonBody(map[string]float64{"amount": 6000}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtido %d: %s", w.Code, w.Body.String())
	}
}

func TestAuctionHandler_PlaceBid_IDInvalido(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{}}
	r := newTestRouter()
	r.POST("/auctions/:id/bids", injetaUserID(1), h.PlaceBid)

	req := httptest.NewRequest(http.MethodPost, "/auctions/abc/bids", jsonBody(map[string]float64{"amount": 100}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtido %d", w.Code)
	}
}

func TestAuctionHandler_PlaceBid_ErroService(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{placeBidErr: errors.New("leilão encerrado")}}
	r := newTestRouter()
	r.POST("/auctions/:id/bids", injetaUserID(1), h.PlaceBid)

	req := httptest.NewRequest(http.MethodPost, "/auctions/1/bids", jsonBody(map[string]float64{"amount": 100}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtido %d", w.Code)
	}
}

func TestAuctionHandler_GetBids_Sucesso(t *testing.T) {
	bids := []models.Bid{{ID: 1, Amount: 7000}, {ID: 2, Amount: 6000}}
	h := &AuctionHandler{Service: &mockAuctionSvc{getBidsResult: bids}}
	r := newTestRouter()
	r.GET("/auctions/:id/bids", h.GetBids)

	req := httptest.NewRequest(http.MethodGet, "/auctions/1/bids", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("esperado 200, obtido %d", w.Code)
	}
}

func TestAuctionHandler_GetBids_IDInvalido(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{}}
	r := newTestRouter()
	r.GET("/auctions/:id/bids", h.GetBids)

	req := httptest.NewRequest(http.MethodGet, "/auctions/abc/bids", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtido %d", w.Code)
	}
}

func TestAuctionHandler_GetBids_Erro(t *testing.T) {
	h := &AuctionHandler{Service: &mockAuctionSvc{getBidsErr: errors.New("db error")}}
	r := newTestRouter()
	r.GET("/auctions/:id/bids", h.GetBids)

	req := httptest.NewRequest(http.MethodGet, "/auctions/1/bids", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtido %d", w.Code)
	}
}
