package handler

import (
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"github.com/go-chi/chi/v5"
)

type OrdersHandler struct {
	service ports.OrderService
	router  *chi.Mux
}

func NewOrdersHandler(service ports.OrderService) *OrdersHandler {
	return &OrdersHandler{
		service: service,
	}
}

func (o OrdersHandler) InitRoute() {
	
}
