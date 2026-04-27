package handler

import (
	"database/sql"

	"go-icarros/internal/middleware"
	"go-icarros/internal/repository"
	"go-icarros/internal/service"
	"go-icarros/internal/ws"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	DB     *sql.DB
	Queue  *service.QueueService
	Logger *service.LogService
	Hub    *ws.Hub
}

func RegisterRoutes(r *gin.Engine, d Deps) {
	// --- usuários ---
	userRepo := &repository.UserRepository{DB: d.DB}
	userSvc := &service.UserService{Repo: userRepo}
	userH := &UserHandler{Service: userSvc}

	// --- carros + leilões ---
	carRepo := &repository.CarRepository{DB: d.DB}
	carSvc := &service.CarService{Repo: carRepo}

	auctionRepo := &repository.AuctionRepository{DB: d.DB}
	bidRepo := &repository.BidRepository{DB: d.DB}
	auctionSvc := &service.AuctionService{
		AuctionRepo: auctionRepo,
		BidRepo:     bidRepo,
		Hub:         d.Hub,
		Logger:      d.Logger,
		Publisher:   d.Queue,
	}

	carH := &CarHandler{
		Service:    carSvc,
		AuctionSvc: auctionSvc,
		Logger:     d.Logger,
		Publisher:  d.Queue,
	}
	auctionH := &AuctionHandler{Service: auctionSvc}

	// --- logs ---
	logRepo := &repository.LogRepository{DB: d.DB}
	logSvc := &service.LogService{Repo: logRepo}
	logH := &LogHandler{Service: logSvc}

	// --- websocket ---
	wsH := &WSHandler{Hub: d.Hub}

	// rotas públicas
	r.POST("/login", userH.Login)
	r.POST("/register", userH.Register)

	// rotas de admin
	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("/users", userH.Create)
		admin.GET("/users", userH.List)
		admin.GET("/users/:id", userH.GetByID)
		admin.PUT("/users/:id", userH.Update)
		admin.DELETE("/users/:id", userH.Delete)

		admin.GET("/logs", logH.List)
	}

	// rotas autenticadas
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", userH.GetMe)
		auth.PUT("/me", userH.UpdateMe)

		auth.POST("/cars", carH.Create)
		auth.GET("/cars", carH.List)
		auth.GET("/cars/my", carH.GetMyCars)
		auth.GET("/cars/:id", carH.GetByID)
		auth.PUT("/cars/:id", carH.Update)
		auth.DELETE("/cars/:id", carH.Delete)

		auth.GET("/auctions", auctionH.List)
		auth.GET("/auctions/:id", auctionH.GetByID)
		auth.POST("/auctions/:id/bids", auctionH.PlaceBid)
		auth.GET("/auctions/:id/bids", auctionH.GetBids)
	}

	// websocket (auth via query param ?token=)
	r.GET("/ws/auctions/:id", wsH.ServeAuction)
}
