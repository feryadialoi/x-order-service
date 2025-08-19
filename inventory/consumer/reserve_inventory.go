package consumer

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/kafka"
)

type ReserveInventory struct {
	producer kafka.Producer[any]
}

func NewReserveInventoryConsumer(producer kafka.Producer[any]) *ReserveInventory {
	return &ReserveInventory{
		producer: producer,
	}
}

func (c *ReserveInventory) Handle(ctx context.Context, key string, reserverInventory contract.ReserveInventory, header kafka.Header) error {
	if header[contract.HeaderEventType] == (contract.ReserveInventory{}).GetEventType() {
		slog.Info("reserve inventory",
			slog.Any("key", key),
			slog.Any("reserver_inventory", reserverInventory),
			slog.Any("header", header),
		)

		if rand.Intn(100) < 95 {
			return c.producer.Send(ctx, reserverInventory.OrderID, contract.InventoryReserved{
				OrderID: reserverInventory.OrderID,
			}, kafka.Header{
				contract.HeaderEventType: contract.InventoryReserved{}.GetEventType(),
			})
		}

		return c.producer.Send(ctx, reserverInventory.OrderID, contract.OrderFailed{
			OrderID: reserverInventory.OrderID,
			Reason:  "Inventory not available",
		}, kafka.Header{
			contract.HeaderEventType: contract.OrderFailed{}.GetEventType(),
		})
	}
	return nil
}
