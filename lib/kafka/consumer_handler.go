package kafka

import (
	"context"
)

type ConsumerHandler[T any] interface {
	Handle(ctx context.Context, key string, value T, header Header) error
}

type ConsumerHandlerFunc[T any] func(ctx context.Context, key string, value T, header Header) error

func (f ConsumerHandlerFunc[T]) Handle(ctx context.Context, key string, value T, header Header) error {
	return f(ctx, key, value, header)
}

type Unmarshaler interface {
	Unmarshal(msg []byte) error
}
