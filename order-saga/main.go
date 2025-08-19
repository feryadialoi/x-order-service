package main

import (
	"database/sql"

	"github.com/feryadialoi/x-order-service/lib/graceful"
	"github.com/feryadialoi/x-order-service/lib/must"
	"github.com/feryadialoi/x-order-service/lib/saga"
	"github.com/feryadialoi/x-order-service/order-saga/entity"
	"github.com/feryadialoi/x-order-service/order-saga/statemachine"
	_ "github.com/lib/pq"
)

func main() {
	db := must.Must(sql.Open("postgres", "host=localhost port=5432 user=postgres password= dbname=postgres sslmode=disable"))

	defer func() { _ = db.Close() }()

	s := saga.NewSaga[entity.OrderState](
		saga.WithDatabase[entity.OrderState](db, "order_states"),
		saga.WithKafka[entity.OrderState]("localhost:9092", "order-saga", "order-saga"),
		saga.WithStateMachine[entity.OrderState](statemachine.NewOrder()),
	)

	graceful.Graceful(map[string]graceful.Process{
		"order-saga": s,
	})
}
