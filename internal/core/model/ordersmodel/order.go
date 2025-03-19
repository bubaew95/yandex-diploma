package ordersmodel

type Order struct {
	Number int64 `db:"number"`
	UserID int64 `db:"user_id"`
}
