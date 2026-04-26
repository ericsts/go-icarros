package repository

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"go-icarros/internal/models"
)

type LogRepository struct {
	DB *sql.DB
}

func (r *LogRepository) Create(entry *models.EventLog) error {
	meta, _ := json.Marshal(entry.Metadata)
	_, err := r.DB.Exec(
		"INSERT INTO event_logs(level, event, message, metadata) VALUES($1,$2,$3,$4)",
		entry.Level, entry.Event, entry.Message, meta,
	)
	return err
}

func (r *LogRepository) FindAll(level, event string, limit int) ([]models.EventLog, error) {
	query := "SELECT id, level, event, message, metadata, created_at FROM event_logs WHERE 1=1"
	args := []any{}

	if level != "" {
		args = append(args, level)
		query += " AND level=$" + strconv.Itoa(len(args))
	}
	if event != "" {
		args = append(args, "%"+event+"%")
		query += " AND event ILIKE $" + strconv.Itoa(len(args))
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		args = append(args, limit)
		query += " LIMIT $" + strconv.Itoa(len(args))
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.EventLog
	for rows.Next() {
		var entry models.EventLog
		var metaRaw []byte
		if err := rows.Scan(&entry.ID, &entry.Level, &entry.Event, &entry.Message, &metaRaw, &entry.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(metaRaw, &entry.Metadata)
		logs = append(logs, entry)
	}
	return logs, nil
}
