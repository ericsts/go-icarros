package service

import (
	"errors"
	"testing"
	"time"

	"go-icarros/internal/models"
)

type mockAuctionRepo struct {
	createErr           error
	findAllResult       []models.Auction
	findAllErr          error
	findByIDResult      *models.Auction
	findByIDErr         error
	findExpiredResult   []models.Auction
	findExpiredErr      error
	updateStatusErr     error
	findOpenByCarResult *models.Auction
	findOpenByCarErr    error
}

func (m *mockAuctionRepo) Create(_ *models.Auction) error {
	return m.createErr
}
func (m *mockAuctionRepo) FindAll() ([]models.Auction, error) {
	return m.findAllResult, m.findAllErr
}
func (m *mockAuctionRepo) FindByID(_ int) (*models.Auction, error) {
	return m.findByIDResult, m.findByIDErr
}
func (m *mockAuctionRepo) FindExpired() ([]models.Auction, error) {
	return m.findExpiredResult, m.findExpiredErr
}
func (m *mockAuctionRepo) UpdateStatus(_ int, _ string) error {
	return m.updateStatusErr
}
func (m *mockAuctionRepo) FindOpenByCarID(_ int) (*models.Auction, error) {
	return m.findOpenByCarResult, m.findOpenByCarErr
}

type mockBidRepo struct {
	createErr             error
	findByAuctionIDResult []models.Bid
	findByAuctionIDErr    error
	findHighestResult     *models.Bid
	findHighestErr        error
}

func (m *mockBidRepo) Create(_ *models.Bid) error { return m.createErr }
func (m *mockBidRepo) FindByAuctionID(_ int) ([]models.Bid, error) {
	return m.findByAuctionIDResult, m.findByAuctionIDErr
}
func (m *mockBidRepo) FindHighestByAuctionID(_ int) (*models.Bid, error) {
	return m.findHighestResult, m.findHighestErr
}

type mockBroadcaster struct {
	lastAuctionID int
}

func (m *mockBroadcaster) Broadcast(auctionID int, _ any) { m.lastAuctionID = auctionID }

type mockLogger struct{}

func (m *mockLogger) Info(_ string, _ string, _ map[string]any)  {}
func (m *mockLogger) Warn(_ string, _ string, _ map[string]any)  {}
func (m *mockLogger) Error(_ string, _ string, _ map[string]any) {}

func newAuctionSvc(aRepo *mockAuctionRepo, bRepo *mockBidRepo) *AuctionService {
	return &AuctionService{
		AuctionRepo: aRepo,
		BidRepo:     bRepo,
		Hub:         &mockBroadcaster{},
		Logger:      &mockLogger{},
	}
}

func TestAuctionService_CreateForCar(t *testing.T) {
	svc := newAuctionSvc(&mockAuctionRepo{}, &mockBidRepo{})
	endsAt := time.Now().Add(24 * time.Hour)

	a, err := svc.CreateForCar(1, endsAt, 5000)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if a.CarID != 1 {
		t.Errorf("esperado CarID=1, obtido %d", a.CarID)
	}
	if a.Status != "open" {
		t.Errorf("esperado Status=open, obtido %s", a.Status)
	}
	if a.MinBid != 5000 {
		t.Errorf("esperado MinBid=5000, obtido %.2f", a.MinBid)
	}
}

func TestAuctionService_CreateForCar_ErroRepo(t *testing.T) {
	svc := newAuctionSvc(&mockAuctionRepo{createErr: errors.New("db error")}, &mockBidRepo{})

	_, err := svc.CreateForCar(1, time.Now().Add(time.Hour), 1000)

	if err == nil {
		t.Fatal("deveria retornar erro")
	}
}

func TestAuctionService_GetAll(t *testing.T) {
	auctions := []models.Auction{{ID: 1}, {ID: 2}}
	svc := newAuctionSvc(&mockAuctionRepo{findAllResult: auctions}, &mockBidRepo{})

	result, err := svc.GetAll()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("esperado 2 leilões, obtido %d", len(result))
	}
}

func TestAuctionService_GetByID(t *testing.T) {
	svc := newAuctionSvc(
		&mockAuctionRepo{findByIDResult: &models.Auction{ID: 3, Status: "open"}},
		&mockBidRepo{},
	)

	a, err := svc.GetByID(3)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if a.ID != 3 {
		t.Errorf("esperado ID=3, obtido %d", a.ID)
	}
}

func TestAuctionService_PlaceBid_Sucesso(t *testing.T) {
	hub := &mockBroadcaster{}
	auction := &models.Auction{
		ID:     1,
		Status: "open",
		EndsAt: time.Now().Add(time.Hour),
		MinBid: 5000,
	}
	svc := &AuctionService{
		AuctionRepo: &mockAuctionRepo{findByIDResult: auction},
		BidRepo:     &mockBidRepo{},
		Hub:         hub,
		Logger:      &mockLogger{},
	}

	bid, err := svc.PlaceBid(1, 2, 6000)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if bid.Amount != 6000 {
		t.Errorf("esperado Amount=6000, obtido %.2f", bid.Amount)
	}
	if bid.UserID != 2 {
		t.Errorf("esperado UserID=2, obtido %d", bid.UserID)
	}
	if hub.lastAuctionID != 1 {
		t.Errorf("hub deveria ter recebido broadcast do leilão 1")
	}
}

func TestAuctionService_PlaceBid_LeilaoNaoEncontrado(t *testing.T) {
	svc := newAuctionSvc(
		&mockAuctionRepo{findByIDErr: errors.New("not found")},
		&mockBidRepo{},
	)

	_, err := svc.PlaceBid(99, 1, 5000)

	if err == nil {
		t.Fatal("deveria retornar erro")
	}
}

func TestAuctionService_PlaceBid_LeilaoEncerrado(t *testing.T) {
	auction := &models.Auction{ID: 1, Status: "closed", EndsAt: time.Now().Add(time.Hour)}
	svc := newAuctionSvc(&mockAuctionRepo{findByIDResult: auction}, &mockBidRepo{})

	_, err := svc.PlaceBid(1, 2, 5000)

	if err == nil || err.Error() != "leilão encerrado" {
		t.Fatalf("esperado 'leilão encerrado', obtido: %v", err)
	}
}

func TestAuctionService_PlaceBid_LeilaoExpirado(t *testing.T) {
	auction := &models.Auction{
		ID:     1,
		Status: "open",
		EndsAt: time.Now().Add(-time.Hour),
		MinBid: 5000,
	}
	svc := newAuctionSvc(&mockAuctionRepo{findByIDResult: auction}, &mockBidRepo{})

	_, err := svc.PlaceBid(1, 2, 6000)

	if err == nil || err.Error() != "leilão expirado" {
		t.Fatalf("esperado 'leilão expirado', obtido: %v", err)
	}
}

func TestAuctionService_PlaceBid_LanceAbaixoMinimo(t *testing.T) {
	auction := &models.Auction{
		ID:     1,
		Status: "open",
		EndsAt: time.Now().Add(time.Hour),
		MinBid: 5000,
	}
	svc := newAuctionSvc(&mockAuctionRepo{findByIDResult: auction}, &mockBidRepo{})

	_, err := svc.PlaceBid(1, 2, 4999)

	if err == nil {
		t.Fatal("deveria rejeitar lance abaixo do mínimo")
	}
}

func TestAuctionService_PlaceBid_LanceAbaixoDoAtual(t *testing.T) {
	auction := &models.Auction{
		ID:         1,
		Status:     "open",
		EndsAt:     time.Now().Add(time.Hour),
		MinBid:     5000,
		CurrentBid: 7000,
	}
	svc := newAuctionSvc(&mockAuctionRepo{findByIDResult: auction}, &mockBidRepo{})

	_, err := svc.PlaceBid(1, 2, 6500)

	if err == nil {
		t.Fatal("deveria rejeitar lance abaixo do lance atual")
	}
}

func TestAuctionService_GetBids(t *testing.T) {
	bids := []models.Bid{{ID: 1, Amount: 6000}, {ID: 2, Amount: 5500}}
	svc := newAuctionSvc(
		&mockAuctionRepo{},
		&mockBidRepo{findByAuctionIDResult: bids},
	)

	result, err := svc.GetBids(1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("esperado 2 lances, obtido %d", len(result))
	}
}
