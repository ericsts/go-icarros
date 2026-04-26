package handler

import (
	"time"

	"go-icarros/internal/models"
	"go-icarros/internal/service"
)

// --- interfaces de serviço consumidas pelos handlers ---

type UserSvc interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type CarSvc interface {
	Create(car *models.Car) error
	GetAll() ([]models.Car, error)
	GetByID(id int) (*models.Car, error)
	GetByUserID(userID int) ([]models.Car, error)
	Update(car *models.Car) error
	Delete(id int) error
}

type AuctionSvc interface {
	CreateForCar(carID int, endsAt time.Time, minBid float64) (*models.Auction, error)
	GetAll() ([]models.Auction, error)
	GetByID(id int) (*models.Auction, error)
	PlaceBid(auctionID, userID int, amount float64) (*models.Bid, error)
	GetBids(auctionID int) ([]models.Bid, error)
}

type LogSvc interface {
	GetAll(level, event string, limit int) ([]models.EventLog, error)
}

// Logger e Publisher são reutilizados do pacote service para evitar duplicação.
type Logger = service.Logger
type Publisher = service.Publisher
