package repository

import (
	"testing"
	"time"

	"go-icarros/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBidRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("INSERT INTO bids").
		WithArgs(1, 2, 6000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(5, now))

	repo := &BidRepository{DB: db}
	b := &models.Bid{AuctionID: 1, UserID: 2, Amount: 6000}

	if err := repo.Create(b); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if b.ID != 5 {
		t.Errorf("esperado ID=5, obtido %d", b.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestBidRepository_FindByAuctionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "auction_id", "user_id", "amount", "created_at"}).
		AddRow(1, 1, 2, 7000.0, now).
		AddRow(2, 1, 3, 6000.0, now)

	mock.ExpectQuery("SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE auction_id").
		WithArgs(1).
		WillReturnRows(rows)

	repo := &BidRepository{DB: db}
	bids, err := repo.FindByAuctionID(1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(bids) != 2 {
		t.Errorf("esperado 2 lances, obtido %d", len(bids))
	}
	if bids[0].Amount != 7000 {
		t.Errorf("esperado Amount=7000, obtido %.2f", bids[0].Amount)
	}
}

func TestBidRepository_FindHighestByAuctionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE auction_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "user_id", "amount", "created_at"}).
			AddRow(1, 1, 2, 8000.0, now))

	repo := &BidRepository{DB: db}
	bid, err := repo.FindHighestByAuctionID(1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if bid == nil {
		t.Fatal("esperado bid, obtido nil")
	}
	if bid.Amount != 8000 {
		t.Errorf("esperado Amount=8000, obtido %.2f", bid.Amount)
	}
}

func TestBidRepository_FindHighestByAuctionID_SemLances(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE auction_id").
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "user_id", "amount", "created_at"}))

	repo := &BidRepository{DB: db}
	bid, err := repo.FindHighestByAuctionID(99)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if bid != nil {
		t.Errorf("esperado nil quando não há lances, obtido bid ID=%d", bid.ID)
	}
}
