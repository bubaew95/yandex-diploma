package ports

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	"github.com/bubaew95/yandex-diploma/internal/core/model/ordersmodel"
)

type OrderService interface {
	AddOrdersNumber(ctx context.Context, number string) error
	OrdersByUserId(ctx context.Context) ([]ordersmodel.Orders, error)
}

type OrderRepository interface {
	AddOrdersNumber(ctx context.Context, order orderentity.Order) error
	OrdersByUserId(ctx context.Context, userId int64) ([]ordersmodel.Orders, error)
}
