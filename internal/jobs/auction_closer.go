package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-icarros/internal/service"
)

// StartAuctionCloser verifica a cada minuto se há leilões expirados e os encerra.
func StartAuctionCloser(
	auctionRepo service.AuctionRepo,
	bidRepo service.BidRepo,
	publisher service.Publisher,
	logger service.Logger,
) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			closeExpired(auctionRepo, bidRepo, publisher, logger)
		}
	}()
	log.Println("auction closer: iniciado, verificando a cada minuto")
}

func closeExpired(
	auctionRepo service.AuctionRepo,
	bidRepo service.BidRepo,
	publisher service.Publisher,
	logger service.Logger,
) {
	expired, err := auctionRepo.FindExpired()
	if err != nil {
		logger.Error("auction.closer", "erro ao buscar leilões expirados", map[string]any{"error": err.Error()})
		return
	}

	for _, a := range expired {
		if err := auctionRepo.UpdateStatus(a.ID, "closed"); err != nil {
			logger.Error("auction.closer", "erro ao encerrar leilão", map[string]any{"auction_id": a.ID, "error": err.Error()})
			continue
		}

		highest, err := bidRepo.FindHighestByAuctionID(a.ID)
		if err != nil || highest == nil {
			logger.Info("auction.closed", fmt.Sprintf("leilão %d encerrado sem lances", a.ID), map[string]any{"auction_id": a.ID})
			continue
		}

		event, _ := json.Marshal(map[string]any{
			"auction_id": a.ID,
			"car_id":     a.CarID,
			"winner_id":  highest.UserID,
			"amount":     highest.Amount,
		})
		publisher.Publish("auction.closed", event)

		logger.Info("auction.closed",
			fmt.Sprintf("leilão %d encerrado — vencedor: usuário %d com R$ %.2f", a.ID, highest.UserID, highest.Amount),
			map[string]any{"auction_id": a.ID, "winner_id": highest.UserID, "amount": highest.Amount},
		)
	}
}
