package usermodel

import "time"

type Withdraw struct {
	UserID      int64     `json:"userid"`
	OrderNumber int64     `db:"order_number"`
	Amount      float64   `db:"amount"`
	ProcessedAt time.Time `db:"processed_at"`
}
