package main

import (
	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/inventory/consumer"
	"github.com/feryadialoi/x-order-service/lib/graceful"
	"github.com/feryadialoi/x-order-service/lib/kafka"
	"github.com/feryadialoi/x-order-service/lib/must"
)

func main() {
	producer := must.Must(kafka.NewSyncProducer[any](
		[]string{"localhost:9092"},
		"order-saga",
	))

	reserveInventory := consumer.NewReserveInventoryConsumer(producer)

	reserveInventoryConsumer := must.Must(kafka.NewConsumer[contract.ReserveInventory](
		[]string{"localhost:9092"},
		"order-saga",
		"reserve-inventory",
		reserveInventory,
		kafka.WithOffset[contract.ReserveInventory](kafka.OffsetOldest),
	))

	graceful.Graceful(map[string]graceful.Process{
		"reserve-inventory-consumer": reserveInventoryConsumer,
	})
}
