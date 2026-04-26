package repository

import (
	"database/sql"

	"go-icarros/internal/models"
)

type AuctionRepository struct {
	DB *sql.DB
}

func (r *AuctionRepository) Create(a *models.Auction) error {
	return r.DB.QueryRow(
		"INSERT INTO auctions(car_id, ends_at, status, min_bid) VALUES($1,$2,$3,$4) RETURNING id, created_at",
		a.CarID, a.EndsAt, a.Status, a.MinBid,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *AuctionRepository) FindAll() ([]models.Auction, error) {
	rows, err := r.DB.Query(`
		SELECT a.id, a.car_id, a.ends_at, a.status, a.min_bid, a.created_at,
		       COALESCE(MAX(b.amount), 0) AS current_bid,
		       COUNT(b.id) AS total_bids
		FROM auctions a
		LEFT JOIN bids b ON b.auction_id = a.id
		GROUP BY a.id
		ORDER BY a.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []models.Auction
	for rows.Next() {
		var a models.Auction
		if err := rows.Scan(&a.ID, &a.CarID, &a.EndsAt, &a.Status, &a.MinBid, &a.CreatedAt, &a.CurrentBid, &a.TotalBids); err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}
	return auctions, nil
}

func (r *AuctionRepository) FindByID(id int) (*models.Auction, error) {
	var a models.Auction
	err := r.DB.QueryRow(`
		SELECT a.id, a.car_id, a.ends_at, a.status, a.min_bid, a.created_at,
		       COALESCE(MAX(b.amount), 0) AS current_bid,
		       COUNT(b.id) AS total_bids
		FROM auctions a
		LEFT JOIN bids b ON b.auction_id = a.id
		WHERE a.id = $1
		GROUP BY a.id
	`, id).Scan(&a.ID, &a.CarID, &a.EndsAt, &a.Status, &a.MinBid, &a.CreatedAt, &a.CurrentBid, &a.TotalBids)
	return &a, err
}

func (r *AuctionRepository) FindExpired() ([]models.Auction, error) {
	rows, err := r.DB.Query(
		"SELECT id, car_id, ends_at, status, min_bid, created_at FROM auctions WHERE status='open' AND ends_at < NOW()",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []models.Auction
	for rows.Next() {
		var a models.Auction
		if err := rows.Scan(&a.ID, &a.CarID, &a.EndsAt, &a.Status, &a.MinBid, &a.CreatedAt); err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}
	return auctions, nil
}

func (r *AuctionRepository) UpdateStatus(id int, status string) error {
	_, err := r.DB.Exec("UPDATE auctions SET status=$1 WHERE id=$2", status, id)
	return err
}
