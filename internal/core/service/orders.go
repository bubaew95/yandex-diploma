package service

import (
	"context"
	"fmt"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/bubaew95/yandex-diploma/internal/core/constants"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/ordersdto"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/systemdto"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/calcsystementity"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/orderentity"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/ordersmodel"
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"go.uber.org/zap"
	"gopkg.in/resty.v1"
	"net/http"
	"regexp"
	"strconv"
	"time"
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

	user, ok := ctx.Value(constants.UserKey).(userentity.User)
	if !ok {
		return apperrors.ErrUserNotFound
	}

	reg := regexp.MustCompile("[^0-9]+")
	num := reg.ReplaceAllString(number, "")

	err := goluhn.Validate(num)
	if err != nil {
		return apperrors.ErrInvalidOrderNumber
	}

	userOrder := ordersmodel.Order{
		UserID: user.ID,
		Number: num,
	}

	return s.repo.AddOrdersNumber(ctx, userOrder)
}

func (s OrdersService) OrdersByUserID(ctx context.Context) ([]ordersdto.Orders, error) {
	user, ok := ctx.Value(constants.UserKey).(userentity.User)
	if !ok {
		return nil, apperrors.ErrUserNotFound
	}

	orders, err := s.repo.OrdersByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s OrdersService) OrdersWithoutAccrual(ctx context.Context) ([]orderentity.OrderDetails, error) {
	return s.repo.OrdersWithoutAccrual(ctx)
}

func (s OrdersService) UpdateOrderByID(ctx context.Context, userID int64, cs systemdto.CalculationSystem) error {
	return s.repo.UpdateOrderByID(ctx, userID, cs)
}

func (s OrdersService) GetPointByNumber(ctx context.Context, number string) (calcsystementity.CalculationSystem, error) {
	calcSystemURL := fmt.Sprintf("%s/api/orders/%s", s.config.AccrualAddress, number)

	var calcResponse systemdto.CalculationSystem

	client := resty.New()
	res, err := client.R().
		SetContext(ctx).
		Get(calcSystemURL)
	if err != nil {
		return calcsystementity.CalculationSystem{}, err
	}

	logger.Log.Info("Система рассчета", zap.Any("calcResponse", res))

	retry := res.Header().Get("Retry-After")
	return calcsystementity.CalculationSystem{
		CalculationSystem: &calcResponse,
		StatusCode:        res.StatusCode(),
		Retry:             retry,
	}, err
}

func (s OrdersService) processOrder(
	ctx context.Context,
	order orderentity.OrderDetails,
	resultCh chan error,
	retryQueueCh chan calcsystementity.RetryOrder) {

	res, err := s.GetPointByNumber(ctx, order.Number)
	if err != nil {
		resultCh <- err
	}

	if res.StatusCode == http.StatusTooManyRequests {
		retrySeconds, err := strconv.ParseInt(res.Retry, 10, 64)
		if err != nil {
			resultCh <- err
		}

		retryAfter := time.Now().Add(time.Duration(retrySeconds) * time.Second)

		select {
		case retryQueueCh <- calcsystementity.RetryOrder{Order: order, RetryTime: retryAfter}:
			logger.Log.Debug("Заказ получил 429, повтор через сек", zap.String("number", order.Number), zap.Int64("time", retrySeconds))
		default:
			logger.Log.Debug("Очередь переполнена! Заказ пропущен", zap.String("number", order.Number))
		}
		return
	}

	err = s.UpdateOrderByID(ctx, order.UserID, *res.CalculationSystem)
	if err != nil {
		resultCh <- err
	}
}

func (s OrdersService) Worker(ctx context.Context, resultCh chan error) {
	orderCh := make(chan orderentity.OrderDetails)
	retryQueueCh := make(chan calcsystementity.RetryOrder, 100)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				orders, err := s.OrdersWithoutAccrual(ctx)
				if err != nil {
					resultCh <- err
					continue
				}
				for _, order := range orders {
					orderCh <- order
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case order := <-orderCh:
				s.processOrder(ctx, order, resultCh, retryQueueCh)

			case retryOrder := <-retryQueueCh:
				time.Sleep(time.Until(retryOrder.RetryTime))
				s.processOrder(ctx, retryOrder.Order, resultCh, retryQueueCh)

			case <-ctx.Done():
				return
			}
		}
	}()
}
