package usermodel

type Balance struct {
	Current   float64 `db:"current"`
	Withdrawn float64 `db:"withdrawn"`
}
