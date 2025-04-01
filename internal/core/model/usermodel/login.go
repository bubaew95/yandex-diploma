package usermodel

type UserLogin struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}
 