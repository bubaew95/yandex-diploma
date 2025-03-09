package orderentity

import "time"

type OrderDetails struct {
	Id        int64
	UserId    int64
	Number    int64
	Status    string
	CreatedAt time.Time
}
