package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/ordersdto"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/ordersmodel"
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"regexp"
	"strconv"
)

type OrdersService struct {
	repo   ports.OrderRepository
	config *conf.Config
}

func NewOrdersService(repo ports.OrderRepository, conf *conf.Config) *OrdersService {
	return &OrdersService{
		repo:   repo,
		config: conf,
	}
}

func (s OrdersService) AddOrdersNumber(ctx context.Context, number string) error {
	if number == "" {
		return apperrors.ErrIncorrectRequest
	}

	user, ok := ctx.Value("user").(userentity.User)
	if !ok {
		return apperrors.ErrUserNotFound
	}

	reg := regexp.MustCompile("[^0-9]+")
	num := reg.ReplaceAllString(number, "")

	err := goluhn.Validate(num)
	if err != nil {
		return apperrors.ErrInvalidOrderNumber
	}

	orderNum, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return err
	}

	userOrder := ordersmodel.Order{
		UserId: user.Id,
		Number: orderNum,
	}

	return s.repo.AddOrdersNumber(ctx, userOrder)
}

func (s OrdersService) OrdersByUserId(ctx context.Context) ([]ordersdto.Orders, error) {
	user, ok := ctx.Value("user").(userentity.User)
	if !ok {
		return nil, apperrors.ErrUserNotFound
	}

	orders, err := s.repo.OrdersByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s OrdersService) OrdersWithoutAccrual(ctx context.Context) ([]orderentity.OrderDetails, error) {
	return s.repo.OrdersWithoutAccrual(ctx)
}

func (s OrdersService) Worker(ctx context.Context, orderId string) {

}
