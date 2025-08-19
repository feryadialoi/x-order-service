package saga

import "time"

type Entity[T any] struct {
	CorrelationID string
	CurrentState  string
	Data          T
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
