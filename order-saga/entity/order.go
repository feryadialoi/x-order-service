package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type OrderState struct {
	OrderID       string          `json:"order_id"`
	OrderTotal    decimal.Decimal `json:"order_total"`
	PaymentID     string          `json:"payment_id"`
	CustomerEmail string          `json:"customer_email"`
	OrderAt       time.Time       `json:"order_at"`
}

// Scan implements Scanner (for reading from DB)
func (p *OrderState) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, p)
}

// Value implements Valuer (for writing to DB)
func (p OrderState) Value() (driver.Value, error) {
	return json.Marshal(p)
}
