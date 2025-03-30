package domain

import (
	"context"
)

type TemperatureRepository interface {
	Add(ctx context.Context, record *DefaultRecord) error
}
