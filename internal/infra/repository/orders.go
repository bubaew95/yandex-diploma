package repository

import (
	"context"
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

func (o OrdersRepository) AddOrdersNumber(ctx context.Context, order orderentity.Order) error {
	findOrder, err := o.GetOrderByNumber(ctx, order.Number)
	if err == nil {
		return checkAddedOrder(findOrder, order.UserId)
	}

	sqlQuery := `INSERT INTO orders (user_id, number, status) VALUES ($1, $2, $3)`
	_, err = o.db.ExecContext(ctx, sqlQuery, order.UserId, order.Number, StatusProcessing)
	if err != nil {
		return err
	}

	return nil
}

func checkAddedOrder(order orderentity.OrderDetails, userId int64) error {
	if order.UserId != userId {
		return apperrors.ErrOrderAddedAnotherUser
	}

	return apperrors.ErrOrderAddedThisUser
}

func (o OrdersRepository) GetOrderByNumber(ctx context.Context, number int64) (orderentity.OrderDetails, error) {
	sqlQuery := `SELECT id, status, user_id, uploaded_at FROM orders WHERE number = $1`
	row := o.db.QueryRowContext(ctx, sqlQuery, number)

	if row.Err() != nil {
		return orderentity.OrderDetails{}, row.Err()
	}

	var order orderentity.OrderDetails
	err := row.Scan(&order.Id, &order.Status, &order.UserId, &order.CreatedAt)
	if err != nil {
		return orderentity.OrderDetails{}, apperrors.ErrOrderNotFound
	}

	return order, nil
}

func (o OrdersRepository) OrdersByUserId(ctx context.Context, userId int64) ([]ordersmodel.Orders, error) {
	sqlQuery := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := o.db.QueryContext(ctx, sqlQuery, userId)
	if err != nil {
		return []ordersmodel.Orders{}, err
	}
	defer rows.Close()

	orders := make([]ordersmodel.Orders, 0)
	for rows.Next() {
		var order ordersmodel.Orders
		var uploadedAt time.Time

		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &uploadedAt)
		if err != nil {
			return []ordersmodel.Orders{}, err
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
