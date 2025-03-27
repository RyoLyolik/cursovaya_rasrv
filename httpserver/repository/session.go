package repository

import (
	"context"
	"httpserver/domain"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	database   *redis.Client
	logger     *slog.Logger
	expiration time.Duration
}

func NewSessionRepository(db *redis.Client, log *slog.Logger, ex time.Duration) domain.SessionRepository {
	return &sessionRepository{
		database:   db,
		logger:     log,
		expiration: ex,
	}
}

func (sr *sessionRepository) Save(ctx context.Context, session_id string, id int) error {
	status := sr.database.Set(ctx, session_id, id, sr.expiration)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (sr *sessionRepository) Get(ctx context.Context, session_id string) (int, error) {
	stringCmd := sr.database.Get(ctx, session_id)
	if stringCmd.Err() != nil {
		return 0, stringCmd.Err()
	}
	id, err := stringCmd.Int()
	if err != nil {
		return 0, err
	}
	return id, nil
}
