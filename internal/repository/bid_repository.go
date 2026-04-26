package repository

import (
	"database/sql"

	"go-icarros/internal/models"
)

type BidRepository struct {
	DB *sql.DB
}

func (r *BidRepository) Create(b *models.Bid) error {
	return r.DB.QueryRow(
		"INSERT INTO bids(auction_id, user_id, amount) VALUES($1,$2,$3) RETURNING id, created_at",
		b.AuctionID, b.UserID, b.Amount,
	).Scan(&b.ID, &b.CreatedAt)
}

func (r *BidRepository) FindByAuctionID(auctionID int) ([]models.Bid, error) {
	rows, err := r.DB.Query(
		"SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE auction_id=$1 ORDER BY amount DESC",
		auctionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var b models.Bid
		if err := rows.Scan(&b.ID, &b.AuctionID, &b.UserID, &b.Amount, &b.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, b)
	}
	return bids, nil
}

func (r *BidRepository) FindHighestByAuctionID(auctionID int) (*models.Bid, error) {
	var b models.Bid
	err := r.DB.QueryRow(
		"SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE auction_id=$1 ORDER BY amount DESC LIMIT 1",
		auctionID,
	).Scan(&b.ID, &b.AuctionID, &b.UserID, &b.Amount, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &b, err
}
