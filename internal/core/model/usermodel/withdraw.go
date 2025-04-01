package usermodel

import "time"

type Withdraw struct {
	UserID      int64     `db:"userid"`
	OrderNumber int64     `db:"order_number"`
	Amount      float64   `db:"amount"`
	ProcessedAt time.Time `db:"processed_at"`
}
