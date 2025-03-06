package repository

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/entity"
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

func (r *UserRepository) Registration(ctx context.Context, u model.SignUp) (entity.User, error) {
	sqlString := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id"

	row := r.db.QueryRowContext(ctx, sqlString, u.Login, u.Password)
	if row.Err() != nil {
		return entity.User{}, u.CheckLogin(row.Err())
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return entity.User{}, err
	}

	return entity.User{
		Id:    id,
		Login: u.Login,
	}, nil
}
