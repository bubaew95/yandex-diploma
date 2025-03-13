package repository

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/responsedto"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
	"github.com/bubaew95/yandex-diploma/internal/infra"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
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
		return userentity.User{}, apperrors.ErrIncorrectUser
	}

	return user, nil
}

func (r UserRepository) GetUserBalance(ctx context.Context, userID int64) (usermodel.Balance, error) {
	sqlQuery := `
		SELECT ub.balance, ub.withdrawn FROM user_balance ub
			INNER JOIN users u ON u.id = ub.user_id
				WHERE ub.user_id = $1
	  `

	row := r.db.QueryRowContext(ctx, sqlQuery, userID)
	if row.Err() != nil {
		return usermodel.Balance{}, row.Err()
	}

	var balance usermodel.Balance
	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return usermodel.Balance{}, err
	}

	return balance, nil
}

func (r UserRepository) BalanceWithdraw(ctx context.Context, ur usermodel.Withdraw) error {
	balance, err := r.GetUserBalance(ctx, ur.UserID)
	if err != nil {
		return err
	}

	if ur.Amount > balance.Current {
		return apperrors.ErrInsufficientErr
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	updateBalanceSql := `
		UPDATE user_balance SET balance = (balance - $1), withdrawn = (withdrawn + $1) 
			WHERE user_id = $2
	`
	_, err = tx.ExecContext(ctx, updateBalanceSql, ur.Amount, ur.UserID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return apperrors.ErrBalanceUpdate
	}

	err = r.AddWithdraw(ctx, ur)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r UserRepository) AddWithdraw(ctx context.Context, ur usermodel.Withdraw) error {
	insertWithdrawSql := `
		INSERT INTO withdraws (user_id, order_number, amount) 
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, insertWithdrawSql, ur.UserID, ur.OrderNumber, ur.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetWithdrawals(ctx context.Context) ([]responsedto.Withdraw, error) {
	sqlString := "SELECT order_number, amount, processed_at FROM withdraws ORDER BY processed_at DESC"
	rows, err := r.db.QueryContext(ctx, sqlString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	withdraws := make([]responsedto.Withdraw, 0)
	for rows.Next() {
		var withdraw responsedto.Withdraw
		var processedAt time.Time
		
		err := rows.Scan(&withdraw.OrderNumber, &withdraw.Sum, &processedAt)
		if err != nil {
			return nil, err
		}

		withdraw.ProcessedAt = processedAt.Format(time.RFC3339)

		withdraws = append(withdraws, withdraw)
	}

	return withdraws, nil
}
