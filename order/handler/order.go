package handler

import (
	"context"

	"github.com/feryadialoi/x-order-service/contract"
	"github.com/feryadialoi/x-order-service/lib/kafka"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SubmitOrderRequest struct {
	Total decimal.Decimal `json:"total"`
	Email string          `json:"email"`
}

type SubmitOrderResponse struct {
	OrderID string `json:"order_id"`
}

type Order struct {
	producer kafka.Producer[contract.OrderSubmitted]
}

func NewOrder(producer kafka.Producer[contract.OrderSubmitted]) *Order {
	return &Order{
		producer: producer,
	}
}

func (o *Order) SubmitOrder(ctx context.Context, req SubmitOrderRequest) (SubmitOrderResponse, error) {
	orderID := uuid.NewString()

	if err := o.producer.Send(ctx, orderID, contract.OrderSubmitted{
		OrderID: orderID,
		Total:   req.Total,
		Email:   req.Email,
	}, kafka.Header{
		contract.HeaderEventType: contract.OrderSubmitted{}.GetEventType(),
	}); err != nil {
		return SubmitOrderResponse{}, err
	}

	response := SubmitOrderResponse{
		OrderID: orderID,
	}

	return response, nil
}
