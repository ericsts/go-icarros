package jobs

import (
	"errors"
	"testing"
	"time"

	"go-icarros/internal/models"
)

type mockAuctionRepo struct {
	findExpiredResult []models.Auction
	findExpiredErr    error
	updateStatusErr   error
}

func (m *mockAuctionRepo) Create(_ *models.Auction) error          { return nil }
func (m *mockAuctionRepo) FindAll() ([]models.Auction, error)      { return nil, nil }
func (m *mockAuctionRepo) FindByID(_ int) (*models.Auction, error) { return nil, nil }
func (m *mockAuctionRepo) FindExpired() ([]models.Auction, error) {
	return m.findExpiredResult, m.findExpiredErr
}
func (m *mockAuctionRepo) UpdateStatus(_ int, _ string) error { return m.updateStatusErr }

type mockBidRepo struct {
	findHighestResult *models.Bid
	findHighestErr    error
}

func (m *mockBidRepo) Create(_ *models.Bid) error                  { return nil }
func (m *mockBidRepo) FindByAuctionID(_ int) ([]models.Bid, error) { return nil, nil }
func (m *mockBidRepo) FindHighestByAuctionID(_ int) (*models.Bid, error) {
	return m.findHighestResult, m.findHighestErr
}

type mockPublisher struct {
	published []string
}

func (m *mockPublisher) Publish(queue string, _ []byte) error {
	m.published = append(m.published, queue)
	return nil
}

type mockLogger struct {
	errors []string
	infos  []string
}

func (m *mockLogger) Info(_ string, msg string, _ map[string]any)  { m.infos = append(m.infos, msg) }
func (m *mockLogger) Warn(_ string, _ string, _ map[string]any)    {}
func (m *mockLogger) Error(_ string, msg string, _ map[string]any) { m.errors = append(m.errors, msg) }

func TestCloseExpired_SemLeiloes(t *testing.T) {
	aRepo := &mockAuctionRepo{findExpiredResult: []models.Auction{}}
	closeExpired(aRepo, &mockBidRepo{}, &mockPublisher{}, &mockLogger{})
	// sem pânico = ok
}

func TestCloseExpired_ErroAoBuscarExpirados(t *testing.T) {
	aRepo := &mockAuctionRepo{findExpiredErr: errors.New("db error")}
	logger := &mockLogger{}

	closeExpired(aRepo, &mockBidRepo{}, &mockPublisher{}, logger)

	if len(logger.errors) == 0 {
		t.Error("esperado erro logado quando FindExpired falha")
	}
}

func TestCloseExpired_ComVencedor(t *testing.T) {
	expired := []models.Auction{
		{ID: 1, CarID: 10, EndsAt: time.Now().Add(-time.Hour), Status: "open"},
	}
	highest := &models.Bid{ID: 1, AuctionID: 1, UserID: 5, Amount: 8000}

	aRepo := &mockAuctionRepo{findExpiredResult: expired}
	bRepo := &mockBidRepo{findHighestResult: highest}
	publisher := &mockPublisher{}
	logger := &mockLogger{}

	closeExpired(aRepo, bRepo, publisher, logger)

	if len(publisher.published) != 1 || publisher.published[0] != "auction.closed" {
		t.Errorf("esperado publicar 'auction.closed', obtido: %v", publisher.published)
	}
	if len(logger.infos) == 0 {
		t.Error("esperado log de info com dados do vencedor")
	}
}

func TestCloseExpired_SemLances(t *testing.T) {
	expired := []models.Auction{
		{ID: 2, CarID: 11, EndsAt: time.Now().Add(-time.Hour), Status: "open"},
	}

	aRepo := &mockAuctionRepo{findExpiredResult: expired}
	bRepo := &mockBidRepo{findHighestResult: nil}
	publisher := &mockPublisher{}

	closeExpired(aRepo, bRepo, publisher, &mockLogger{})

	if len(publisher.published) != 0 {
		t.Errorf("não deveria publicar evento sem vencedor, obtido: %v", publisher.published)
	}
}

func TestCloseExpired_ErroAoEncerrarLeilao(t *testing.T) {
	expired := []models.Auction{
		{ID: 3, CarID: 12, Status: "open"},
	}

	aRepo := &mockAuctionRepo{
		findExpiredResult: expired,
		updateStatusErr:   errors.New("db error"),
	}
	logger := &mockLogger{}

	closeExpired(aRepo, &mockBidRepo{}, &mockPublisher{}, logger)

	if len(logger.errors) == 0 {
		t.Error("esperado erro logado ao falhar encerramento do leilão")
	}
}
