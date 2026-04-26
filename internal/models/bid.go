package models

import "time"

type Bid struct {
	ID        int       `json:"id"`
	AuctionID int       `json:"auction_id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
