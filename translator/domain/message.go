package domain

import "time"

type RecordType string

const (
	Temperature RecordType = "temperature"
	Pressure    RecordType = "pressure"
	Flap        RecordType = "flap"
)

type IncomingMessage struct {
	Timestamp  time.Time  `json:"timestamp"`
	RecordID   float32    `json:"record_id"`
	RecordType RecordType `json:"record_type"`
	PositionID int        `json:"pos_id"`
	Value      float64    `json:"value"`
}
