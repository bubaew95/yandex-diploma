package model

import (
	"errors"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type SignUp struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (u SignUp) CheckLogin(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return apperrors.LoginAlreadyExists
	}

	return err
}
