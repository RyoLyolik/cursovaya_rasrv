package domain

import (
	"context"
)

type FlapRepository interface {
	Add(ctx context.Context, record *DefaultRecord) error
}
