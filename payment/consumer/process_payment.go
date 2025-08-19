package consumer

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/kafka"
	"github.com/google/uuid"
)

type ProcessPayment struct {
	producer kafka.Producer[any]
}

func NewProcessPaymentConsumer(producer kafka.Producer[any]) *ProcessPayment {
	return &ProcessPayment{
		producer: producer,
	}
}

func (c *ProcessPayment) Handle(ctx context.Context, key string, processPayment contract.ProcessPayment, header kafka.Header) error {
	if header[contract.HeaderEventType] == (contract.ProcessPayment{}).GetEventType() {
		slog.Info("process payment",
			slog.Any("key", key),
			slog.Any("process_payment", processPayment),
			slog.Any("header", header),
		)

		if rand.Intn(100) < 95 {
			return c.producer.Send(ctx, processPayment.OrderID, contract.PaymentProcessed{
				OrderID:   processPayment.OrderID,
				PaymentID: uuid.NewString(),
			}, kafka.Header{
				contract.HeaderEventType: contract.PaymentProcessed{}.GetEventType(),
			})
		}

		return c.producer.Send(ctx, processPayment.OrderID, contract.OrderFailed{
			OrderID: processPayment.OrderID,
			Reason:  "Payment failed",
		}, kafka.Header{
			contract.HeaderEventType: contract.OrderFailed{}.GetEventType(),
		})
	}
	return nil
}
