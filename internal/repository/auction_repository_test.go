package repository

import (
	"testing"
	"time"

	"go-icarros/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuctionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	endsAt := now.Add(24 * time.Hour)

	mock.ExpectQuery("INSERT INTO auctions").
		WithArgs(1, endsAt, "open", 5000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, now))

	repo := &AuctionRepository{DB: db}
	a := &models.Auction{CarID: 1, EndsAt: endsAt, Status: "open", MinBid: 5000}

	if err := repo.Create(a); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if a.ID != 10 {
		t.Errorf("esperado ID=10, obtido %d", a.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestAuctionRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "car_id", "ends_at", "status", "min_bid", "created_at", "current_bid", "total_bids"}).
		AddRow(1, 1, now, "open", 5000.0, now, 6000.0, 3).
		AddRow(2, 2, now, "closed", 10000.0, now, 0.0, 0)

	mock.ExpectQuery("SELECT a.id").WillReturnRows(rows)

	repo := &AuctionRepository{DB: db}
	auctions, err := repo.FindAll()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(auctions) != 2 {
		t.Errorf("esperado 2 leilões, obtido %d", len(auctions))
	}
	if auctions[0].CurrentBid != 6000 {
		t.Errorf("esperado CurrentBid=6000, obtido %.2f", auctions[0].CurrentBid)
	}
}

func TestAuctionRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	mock.ExpectQuery("SELECT a.id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "car_id", "ends_at", "status", "min_bid", "created_at", "current_bid", "total_bids"}).
			AddRow(1, 2, now, "open", 5000.0, now, 0.0, 0))

	repo := &AuctionRepository{DB: db}
	a, err := repo.FindByID(1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if a.ID != 1 {
		t.Errorf("esperado ID=1, obtido %d", a.ID)
	}
	if a.Status != "open" {
		t.Errorf("esperado Status=open, obtido %s", a.Status)
	}
}

func TestAuctionRepository_FindExpired(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "car_id", "ends_at", "status", "min_bid", "created_at"}).
		AddRow(1, 1, now.Add(-time.Hour), "open", 5000.0, now).
		AddRow(2, 3, now.Add(-2*time.Hour), "open", 8000.0, now)

	mock.ExpectQuery("SELECT id, car_id").WillReturnRows(rows)

	repo := &AuctionRepository{DB: db}
	expired, err := repo.FindExpired()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(expired) != 2 {
		t.Errorf("esperado 2 leilões expirados, obtido %d", len(expired))
	}
}

func TestAuctionRepository_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE auctions SET").
		WithArgs("closed", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &AuctionRepository{DB: db}
	if err := repo.UpdateStatus(1, "closed"); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}
