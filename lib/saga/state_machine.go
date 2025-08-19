package saga

import (
	"context"
)

type StateMachine[T any] interface {
	EventConfigurer() EventConfigurer[T]
	EventActivities() EventActivities[T]
}

type State string

const (
	StateInitially State = "Initially"
)

type Event string

type Context[T any] struct {
	ctx     context.Context
	Entity  Entity[T]
	Message []byte
}

func (c *Context[T]) Context() context.Context {
	return c.ctx
}

type EventConfigurer[T any] map[Event]func(c *Context[T])

type Activity[T any] struct {
	Then         func(c *Context[T]) error
	Publish      func(c *Context[T]) any
	TransitionTo State
	Finalize     bool
}

type EventActivities[T any] map[State]map[Event]Activity[T]
