package repository

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/model"
	"github.com/bubaew95/yandex-diploma/internal/infra"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepository struct {
	db *infra.DataBase
}

func NewUserRepository(db *infra.DataBase) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Registration(ctx context.Context, u model.SignUp) error {
	sqlString := "INSERT INTO users (login, password) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, sqlString, u.Login, u.Password)
	if err != nil {
		return u.CheckLogin(err)
	}

	return nil
}
