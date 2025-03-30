package messagehandler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"translator/domain"
	"translator/repository"
)

type messageHandler struct {
	tr  domain.TemperatureRepository
	pr  domain.PressureRepository
	fr  domain.FlapRepository
	log *slog.Logger
}

func NewMessageHandler(log *slog.Logger, db *sql.DB) *messageHandler {
	return &messageHandler{
		tr:  repository.NewTemperatureRepository(db),
		pr:  repository.NewPressureRepository(db),
		fr:  repository.NewFlapRepository(db),
		log: log,
	}
}

func (mh *messageHandler) HandleMessage(ctx context.Context, message []byte) error {
	var msg domain.IncomingMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}
	record := domain.DefaultRecord{
		Timestamp: msg.Timestamp,
		Position:  msg.PositionID,
		Value:     msg.Value,
	}
	switch msg.RecordType {
	case domain.Temperature:
		if err := mh.tr.Add(ctx, &record); err != nil {
			return err
		}
	case domain.Pressure:
		if err := mh.pr.Add(ctx, &record); err != nil {
			return err
		}
	case domain.Flap:
		if err := mh.fr.Add(ctx, &record); err != nil {
			return err
		}
	default:
		return fmt.Errorf("failed to define record type")
	}
	return nil
}
