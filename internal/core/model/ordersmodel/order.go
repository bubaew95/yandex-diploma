package ordersmodel

type Order struct {
	Number int64 `db:"number"`
	UserId int64 `db:"user_id"`
}
