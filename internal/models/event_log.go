package models

import "time"

type EventLog struct {
	ID        int            `json:"id"`
	Level     string         `json:"level"`
	Event     string         `json:"event"`
	Message   string         `json:"message"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}
