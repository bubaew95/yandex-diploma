package ports

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/ordersdto"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	"github.com/bubaew95/yandex-diploma/internal/core/model/ordersmodel"
)

type OrderService interface {
	AddOrdersNumber(ctx context.Context, number string) error
	OrdersByUserId(ctx context.Context) ([]ordersdto.Orders, error)
	OrdersWithoutAccrual(ctx context.Context) ([]orderentity.OrderDetails, error)
}

type OrderRepository interface {
	AddOrdersNumber(ctx context.Context, order ordersmodel.Order) error
	OrdersByUserId(ctx context.Context, userId int64) ([]ordersdto.Orders, error)
	OrdersWithoutAccrual(ctx context.Context) ([]orderentity.OrderDetails, error)
}
