package service

import "go-icarros/internal/models"

// --- repositórios de usuário e carro ---

type UserRepo interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type CarRepo interface {
	Create(car *models.Car) error
	FindAll() ([]models.Car, error)
	FindByID(id int) (*models.Car, error)
	FindByUserID(userID int) ([]models.Car, error)
	Update(car *models.Car) error
	Delete(id int) error
}

// --- repositórios de leilão ---

type AuctionRepo interface {
	Create(a *models.Auction) error
	FindAll() ([]models.Auction, error)
	FindByID(id int) (*models.Auction, error)
	FindExpired() ([]models.Auction, error)
	UpdateStatus(id int, status string) error
}

type BidRepo interface {
	Create(b *models.Bid) error
	FindByAuctionID(auctionID int) ([]models.Bid, error)
	FindHighestByAuctionID(auctionID int) (*models.Bid, error)
}

type LogRepo interface {
	Create(entry *models.EventLog) error
	FindAll(level, event string, limit int) ([]models.EventLog, error)
}

// --- interfaces transversais (usadas por services e jobs) ---

// Logger grava eventos estruturados no banco e no stdout.
type Logger interface {
	Info(event, message string, metadata map[string]any)
	Warn(event, message string, metadata map[string]any)
	Error(event, message string, metadata map[string]any)
}

// Publisher publica mensagens em filas do RabbitMQ.
type Publisher interface {
	Publish(queue string, body []byte) error
}

// Broadcaster envia dados em tempo real via WebSocket para uma sala de leilão.
type Broadcaster interface {
	Broadcast(auctionID int, data any)
}
