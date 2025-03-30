package domain

import "time"

type DefaultRecord struct {
	Timestamp time.Time
	Position  int
	Value     float64
}
