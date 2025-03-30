package domain

import (
	"context"
)

type PressureRepository interface {
	Add(ctx context.Context, record *DefaultRecord) error
}
