package service

import "github.com/bubaew95/yandex-diploma/internal/core/ports"

type OrdersService struct {
	repo ports.OrderRepository
}

func NewOrdersService(repo ports.OrderRepository) *OrdersService {
	return &OrdersService{repo: repo}
}
