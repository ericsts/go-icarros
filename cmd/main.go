package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go-icarros/internal/db"
	"go-icarros/internal/handler"
	"go-icarros/internal/jobs"
	"go-icarros/internal/repository"
	"go-icarros/internal/service"
	"go-icarros/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.Connect()

	hub := ws.NewHub()
	go hub.Run()

	queue, err := service.NewQueueService(getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"))
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}
	defer queue.Close()

	logRepo := &repository.LogRepository{DB: database}
	logger := &service.LogService{Repo: logRepo}

	email := service.NewEmailService()

	startConsumer(queue, email, database, logger)

	jobs.StartAuctionCloser(
		&repository.AuctionRepository{DB: database},
		&repository.BidRepository{DB: database},
		queue,
		logger,
	)

	r := gin.Default()
	handler.RegisterRoutes(r, handler.Deps{
		DB:     database,
		Queue:  queue,
		Logger: logger,
		Hub:    hub,
	})

	log.Println("servidor rodando em :8080")
	r.Run(":8080")
}

func startConsumer(queue *service.QueueService, email *service.EmailService, database *sql.DB, logger *service.LogService) {
	adminEmail := getEnv("ADMIN_EMAIL", "admin@icarros.com")
	userRepo := &repository.UserRepository{DB: database}

	queue.Consume("car.created", func(body []byte) {
		var ev map[string]any
		json.Unmarshal(body, &ev)
		msg := fmt.Sprintf("Marca: %v | Modelo: %v | Usuário ID: %v", ev["marca"], ev["modelo"], ev["user_id"])
		if err := email.Send(adminEmail, "Novo carro cadastrado", msg); err != nil {
			logger.Error("notification.email", "falha ao enviar email (car.created)", map[string]any{"error": err.Error()})
		} else {
			logger.Info("notification.email", "email enviado ao admin: novo carro", ev)
		}
	})

	queue.Consume("car.deleted", func(body []byte) {
		var ev map[string]any
		json.Unmarshal(body, &ev)
		msg := fmt.Sprintf("O carro ID %v foi removido do sistema.", ev["car_id"])
		if err := email.Send(adminEmail, "Carro removido", msg); err != nil {
			logger.Error("notification.email", "falha ao enviar email (car.deleted)", map[string]any{"error": err.Error()})
		} else {
			logger.Info("notification.email", "email enviado ao admin: carro removido", ev)
		}
	})

	queue.Consume("auction.closed", func(body []byte) {
		var ev map[string]any
		json.Unmarshal(body, &ev)
		winnerID := int(ev["winner_id"].(float64))
		winner, err := userRepo.FindByID(winnerID)
		if err != nil {
			logger.Error("notification.email", "vencedor não encontrado", ev)
			return
		}
		msg := fmt.Sprintf(
			"Olá, %s!\n\nVocê venceu o leilão #%v com o lance de R$ %.2f.\nEntre em contato para combinar a retirada do veículo.",
			winner.Name, ev["auction_id"], ev["amount"],
		)
		if err := email.Send(winner.Email, "Parabéns! Você venceu o leilão!", msg); err != nil {
			logger.Error("notification.email", "falha ao enviar email ao vencedor", map[string]any{"error": err.Error(), "winner_id": winnerID})
		} else {
			logger.Info("notification.email", fmt.Sprintf("email enviado ao vencedor %s", winner.Email), ev)
		}
	})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
