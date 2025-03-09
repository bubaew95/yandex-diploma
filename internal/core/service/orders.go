package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/bubaew95/yandex-diploma/conf"
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
	user, ok := ctx.Value("user").(userentity.User)
	if !ok {
		return apperrors.UserNotFoundErr
	}

	reg := regexp.MustCompile("[^0-9]+")
	num := reg.ReplaceAllString(number, "")

	err := goluhn.Validate(num)
	if err != nil {
		return err
	}

	orderNum, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return err
	}

	userOrder := orderentity.Order{
		UserId: user.Id,
		Number: orderNum,
	}

	return s.repo.AddOrdersNumber(ctx, userOrder)
}

func (s OrdersService) OrdersByUserId(ctx context.Context) ([]ordersmodel.Orders, error) {
	user, ok := ctx.Value("user").(userentity.User)
	if !ok {
		return nil, apperrors.UserNotFoundErr
	}

	orders, err := s.repo.OrdersByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
