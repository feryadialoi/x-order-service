package main

import (
	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/graceful"
	"github.com/feryadialoi/x-order-service/lib/kafka"
	"github.com/feryadialoi/x-order-service/lib/must"
	"github.com/feryadialoi/x-order-service/payment/consumer"
)

func main() {
	producer := must.Must(kafka.NewSyncProducer[any](
		[]string{"localhost:9092"},
		"order-saga",
	))

	processPayment := consumer.NewProcessPaymentConsumer(producer)

	processPaymentConsumer := must.Must(kafka.NewConsumer[contract.ProcessPayment](
		[]string{"localhost:9092"},
		"order-saga",
		"process-payment",
		processPayment,
		kafka.WithOffset[contract.ProcessPayment](kafka.OffsetOldest),
	))

	graceful.Graceful(map[string]graceful.Process{
		"process-payment-consumer": processPaymentConsumer,
	})
}
