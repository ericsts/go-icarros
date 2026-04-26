package models

import "time"

type Auction struct {
	ID         int       `json:"id"`
	CarID      int       `json:"car_id"`
	EndsAt     time.Time `json:"ends_at"`
	Status     string    `json:"status"`
	MinBid     float64   `json:"min_bid"`
	CreatedAt  time.Time `json:"created_at"`
	CurrentBid float64   `json:"current_bid,omitempty"`
	TotalBids  int       `json:"total_bids,omitempty"`
}
