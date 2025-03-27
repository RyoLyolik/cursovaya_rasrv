package domain

import "context"

type SessionRepository interface {
	Save(ctx context.Context, session_id string, id int) error
	Get(ctx context.Context, session_id string) (int, error)
}
