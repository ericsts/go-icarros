package service

import (
	"log"

	"go-icarros/internal/models"
)

type LogService struct {
	Repo LogRepo
}

func (s *LogService) Info(event, message string, metadata map[string]any) {
	s.write("info", event, message, metadata)
}

func (s *LogService) Warn(event, message string, metadata map[string]any) {
	s.write("warn", event, message, metadata)
}

func (s *LogService) Error(event, message string, metadata map[string]any) {
	s.write("error", event, message, metadata)
}

func (s *LogService) GetAll(level, event string, limit int) ([]models.EventLog, error) {
	return s.Repo.FindAll(level, event, limit)
}

func (s *LogService) write(level, event, message string, metadata map[string]any) {
	log.Printf("[%s] %s — %s", level, event, message)
	s.Repo.Create(&models.EventLog{
		Level:    level,
		Event:    event,
		Message:  message,
		Metadata: metadata,
	})
}
