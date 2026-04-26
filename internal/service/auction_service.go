package service

import (
	"errors"
	"fmt"
	"time"

	"go-icarros/internal/models"
)

type AuctionService struct {
	AuctionRepo AuctionRepo
	BidRepo     BidRepo
	Hub         Broadcaster
	Logger      Logger
	Publisher   Publisher
}

func (s *AuctionService) CreateForCar(carID int, endsAt time.Time, minBid float64) (*models.Auction, error) {
	a := &models.Auction{
		CarID:  carID,
		EndsAt: endsAt,
		Status: "open",
		MinBid: minBid,
	}
	if err := s.AuctionRepo.Create(a); err != nil {
		return nil, err
	}
	s.Logger.Info("auction.created", fmt.Sprintf("leilão criado para o carro %d", carID), map[string]any{
		"auction_id": a.ID,
		"car_id":     carID,
		"ends_at":    endsAt,
	})
	return a, nil
}

func (s *AuctionService) GetAll() ([]models.Auction, error) {
	return s.AuctionRepo.FindAll()
}

func (s *AuctionService) GetByID(id int) (*models.Auction, error) {
	return s.AuctionRepo.FindByID(id)
}

func (s *AuctionService) PlaceBid(auctionID, userID int, amount float64) (*models.Bid, error) {
	auction, err := s.AuctionRepo.FindByID(auctionID)
	if err != nil {
		return nil, errors.New("leilão não encontrado")
	}
	if auction.Status != "open" {
		return nil, errors.New("leilão encerrado")
	}
	if time.Now().After(auction.EndsAt) {
		return nil, errors.New("leilão expirado")
	}

	minAmount := auction.MinBid
	if auction.CurrentBid > minAmount {
		minAmount = auction.CurrentBid
	}
	if amount <= minAmount {
		return nil, fmt.Errorf("lance deve ser maior que R$ %.2f", minAmount)
	}

	bid := &models.Bid{AuctionID: auctionID, UserID: userID, Amount: amount}
	if err := s.BidRepo.Create(bid); err != nil {
		return nil, err
	}

	if s.Hub != nil {
		s.Hub.Broadcast(auctionID, bid)
	}

	s.Logger.Info("bid.placed", fmt.Sprintf("lance de R$ %.2f no leilão %d", amount, auctionID), map[string]any{
		"auction_id": auctionID,
		"user_id":    userID,
		"amount":     amount,
	})

	return bid, nil
}

func (s *AuctionService) GetBids(auctionID int) ([]models.Bid, error) {
	return s.BidRepo.FindByAuctionID(auctionID)
}
