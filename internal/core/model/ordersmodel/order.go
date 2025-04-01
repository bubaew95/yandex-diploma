package ordersmodel

type Order struct {
	Number string `db:"number"`
	UserID int64  `db:"user_id"`
}
