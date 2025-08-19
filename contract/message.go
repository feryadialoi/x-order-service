package contract

import "github.com/shopspring/decimal"

const (
	HeaderEventType = "event_type"
)

type Event interface {
	GetEventType() string
}

type OrderSubmitted struct {
	OrderID string          `json:"order_id"`
	Total   decimal.Decimal `json:"total"`
	Email   string          `json:"email"`
}

func (o OrderSubmitted) GetEventType() string { return "OrderSubmitted" }

type ProcessPayment struct {
	OrderID string          `json:"order_id"`
	Amount  decimal.Decimal `json:"amount"`
}

func (o ProcessPayment) GetEventType() string { return "ProcessPayment" }

type PaymentProcessed struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
}

func (o PaymentProcessed) GetEventType() string { return "PaymentProcessed" }

type ReserveInventory struct {
	OrderID string `json:"order_id"`
}

func (o ReserveInventory) GetEventType() string { return "ReserveInventory" }

type InventoryReserved struct {
	OrderID string `json:"order_id"`
}

func (o InventoryReserved) GetEventType() string { return "InventoryReserved" }

type RefundPayment struct {
	OrderID string          `json:"order_id"`
	Amount  decimal.Decimal `json:"amount"`
}

func (o RefundPayment) GetEventType() string { return "RefundPayment" }

type OrderConfirmed struct {
	OrderID string
}

func (o OrderConfirmed) GetEventType() string { return "OrderConfirmed" }

type OrderFailed struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
}

func (o OrderFailed) GetEventType() string { return "OrderFailed" }
