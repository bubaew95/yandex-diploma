package orderentity

import "time"

type OrderDetails struct {
	ID        int64
	UserID    int64
	Number    string
	Status    string
	CreatedAt time.Time
}
