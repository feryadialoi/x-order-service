package statemachine

import (
	"time"

	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/jsonmapper"
	"github.com/feryadialoi/x-order-service/lib/saga"
	"github.com/feryadialoi/x-order-service/order-saga/entity"
)

const (
	ProcessingPayment  saga.State = "ProcessingPayment"
	ReservingInventory saga.State = "ReservingInventory"
	Completed          saga.State = "Completed"
	Failed             saga.State = "Failed"
)

const (
	OrderSubmitted    saga.Event = "OrderSubmitted"
	PaymentProcessed  saga.Event = "PaymentProcessed"
	InventoryReserved saga.Event = "InventoryReserved"
	OrderFailed       saga.Event = "OrderFailed"
)

type Order struct {
	eventConfigurer saga.EventConfigurers[entity.OrderState]
	eventActivities saga.EventActivities[entity.OrderState]
}

func NewOrder() *Order {
	return &Order{
		eventConfigurer: saga.EventConfigurers[entity.OrderState]{
			OrderSubmitted: func(c *saga.Context[entity.OrderState]) {
				c.Entity.CorrelationID = jsonmapper.FromJSON[contract.OrderSubmitted](c.Message).OrderID
			},
			PaymentProcessed: func(c *saga.Context[entity.OrderState]) {
				c.Entity.CorrelationID = jsonmapper.FromJSON[contract.PaymentProcessed](c.Message).OrderID
			},
			InventoryReserved: func(c *saga.Context[entity.OrderState]) {
				c.Entity.CorrelationID = jsonmapper.FromJSON[contract.InventoryReserved](c.Message).OrderID
			},
			OrderFailed: func(c *saga.Context[entity.OrderState]) {
				c.Entity.CorrelationID = jsonmapper.FromJSON[contract.OrderFailed](c.Message).OrderID
			},
		},
		eventActivities: saga.EventActivities[entity.OrderState]{
			saga.StateInitially: {
				OrderSubmitted: {
					Then: func(c *saga.Context[entity.OrderState]) error {
						orderSubmitted := jsonmapper.FromJSON[contract.OrderSubmitted](c.Message)
						c.Entity.Data.OrderID = orderSubmitted.OrderID
						c.Entity.Data.OrderTotal = orderSubmitted.Total
						c.Entity.Data.CustomerEmail = orderSubmitted.Email
						c.Entity.Data.OrderAt = time.Now().UTC()
						return nil
					},
					Publish: func(c *saga.Context[entity.OrderState]) any {
						return contract.ProcessPayment{
							OrderID: c.Entity.CorrelationID,
							Amount:  c.Entity.Data.OrderTotal,
						}
					},
					TransitionTo: ProcessingPayment,
				},
			},
			ProcessingPayment: {
				PaymentProcessed: {
					Publish: func(c *saga.Context[entity.OrderState]) any {
						return contract.ReserveInventory{
							OrderID: c.Entity.CorrelationID,
						}
					},
					TransitionTo: ReservingInventory,
				},
				OrderFailed: {
					TransitionTo: Failed,
					Finalize:     true,
				},
			},
			ReservingInventory: {
				InventoryReserved: {
					Publish: func(c *saga.Context[entity.OrderState]) any {
						return contract.OrderConfirmed{
							OrderID: c.Entity.CorrelationID,
						}
					},
					TransitionTo: Completed,
					Finalize:     true,
				},
				OrderFailed: {
					Publish: func(c *saga.Context[entity.OrderState]) any {
						return contract.RefundPayment{
							OrderID: c.Entity.CorrelationID,
							Amount:  c.Entity.Data.OrderTotal,
						}
					},
					TransitionTo: Failed,
					Finalize:     true,
				},
			},
		},
	}
}

func (o *Order) EventConfigurers() saga.EventConfigurers[entity.OrderState] {
	return o.eventConfigurer
}

func (o *Order) EventActivities() saga.EventActivities[entity.OrderState] {
	return o.eventActivities
}
