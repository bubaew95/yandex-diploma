package handler

import (
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"io"
	"net/http"
)

type OrdersHandler struct {
	service ports.OrderService
}

func NewOrdersHandler(service ports.OrderService) *OrdersHandler {
	return &OrdersHandler{
		service: service,
	}
}

func (o OrdersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, response.Response{
			Status:  "failed",
			Message: err.Error(),
		})
	}
	defer r.Body.Close()

	err = o.service.AddOrdersNumber(r.Context(), string(resp))
	if err != nil {
		HandleErrors(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, response.Response{
		Status:  "success",
		Message: string(resp),
	})
}

func (o OrdersHandler) UserOrders(w http.ResponseWriter, r *http.Request) {
	order, err := o.service.OrdersByUserId(r.Context())
	if err != nil {
		HandleErrors(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, order)
}
