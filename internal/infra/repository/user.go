package repository

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
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

func (r UserRepository) AddUser(ctx context.Context, u usermodel.UserRegistration) (userentity.User, error) {
	sqlString := "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id"

	row := r.db.QueryRowContext(ctx, sqlString, u.Login, u.Password)
	if row.Err() != nil {
		return userentity.User{}, u.CheckLogin(row.Err())
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return userentity.User{}, err
	}

	return userentity.User{
		Id:    id,
		Login: u.Login,
	}, nil
}

func (r UserRepository) FindUserByLoginAndPassword(ctx context.Context, u usermodel.UserLogin) (userentity.User, error) {
	sqlString := "SELECT id, login FROM users WHERE login = $1 AND password = $2"

	row := r.db.QueryRowContext(ctx, sqlString, u.Login, u.Password)
	if row.Err() != nil {
		return userentity.User{}, row.Err()
	}

	var user userentity.User
	if err := row.Scan(&user.Id, &user.Login); err != nil {
		return userentity.User{}, apperrors.UserNotFound
	}

	return user, nil
}
