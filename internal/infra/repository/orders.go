package repository

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/ordersdto"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/systemdto"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/ordersmodel"
	"github.com/bubaew95/yandex-diploma/internal/infra"
	"time"
)

type Status string

var (
	StatusNew        Status = "NEW"
	StatusRegistered Status = "REGISTERED"
	StatusInvalid    Status = "INVALID"
	StatusProcessing Status = "PROCESSING"
	StatusProcess    Status = "PROCESSED"
)

type OrdersRepository struct {
	db *infra.DataBase
}

func NewOrdersRepository(db *infra.DataBase) *OrdersRepository {
	return &OrdersRepository{
		db: db,
	}
}

func (o OrdersRepository) AddOrdersNumber(ctx context.Context, order ordersmodel.Order) error {
	findOrder, err := o.GetOrderByNumber(ctx, order.Number)
	if err == nil {
		return checkAddedOrder(findOrder, order.UserID)
	}

	sqlQuery := `INSERT INTO orders (user_id, number, status) VALUES ($1, $2, $3)`
	_, err = o.db.ExecContext(ctx, sqlQuery, order.UserID, order.Number, StatusNew)
	if err != nil {
		return err
	}

	return nil
}

func checkAddedOrder(order orderentity.OrderDetails, userID int64) error {
	if order.UserID != userID {
		return apperrors.ErrOrderAddedAnotherUser
	}

	return apperrors.ErrOrderAddedThisUser
}

func (o OrdersRepository) GetOrderByNumber(ctx context.Context, number string) (orderentity.OrderDetails, error) {
	sqlQuery := `SELECT id, status, user_id, uploaded_at FROM orders WHERE number = $1`
	row := o.db.QueryRowContext(ctx, sqlQuery, number)

	if row.Err() != nil {
		return orderentity.OrderDetails{}, row.Err()
	}

	var order orderentity.OrderDetails
	err := row.Scan(&order.ID, &order.Status, &order.UserID, &order.CreatedAt)
	if err != nil {
		return orderentity.OrderDetails{}, apperrors.ErrOrderNotFound
	}

	return order, nil
}

func (o OrdersRepository) OrdersByUserID(ctx context.Context, userID int64) ([]ordersdto.Orders, error) {
	sqlQuery := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := o.db.QueryContext(ctx, sqlQuery, userID)
	if err != nil {
		return []ordersdto.Orders{}, err
	}
	defer rows.Close()

	orders := make([]ordersdto.Orders, 0)
	for rows.Next() {
		var order ordersdto.Orders
		var uploadedAt time.Time

		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &uploadedAt)
		if err != nil {
			return []ordersdto.Orders{}, err
		}

		order.UploadedAt = uploadedAt.Format(time.RFC3339)

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return orders, apperrors.ErrOrdersEmpty
	}

	return orders, nil
}

func (o OrdersRepository) OrdersWithoutAccrual(ctx context.Context) ([]orderentity.OrderDetails, error) {
	sqlQuery := `SELECT id, user_id, number FROM orders WHERE status NOT IN ('INVALID', 'PROCESSED') AND accrual = 0 ORDER BY id ASC`

	rows, err := o.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return []orderentity.OrderDetails{}, err
	}
	defer rows.Close()

	orders := make([]orderentity.OrderDetails, 0)
	for rows.Next() {
		var order orderentity.OrderDetails
		if err = rows.Scan(&order.ID, &order.UserID, &order.Number); err != nil {
			return []orderentity.OrderDetails{}, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o OrdersRepository) UpdateOrderByID(ctx context.Context, userID int64, cs systemdto.CalculationSystem) error {
	sqlQuery := `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`

	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sqlQuery, cs.Status, cs.Accrual, cs.Order)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlUpdateUserBalance := `UPDATE user_balance SET balance = balance + $1 WHERE user_id = $2`
	_, err = tx.ExecContext(ctx, sqlUpdateUserBalance, cs.Accrual, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
