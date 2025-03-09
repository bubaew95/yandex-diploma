package handler

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"io"
	"net/http"
	"regexp"
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

	reg := regexp.MustCompile("[^0-9]+")
	orderNum := reg.ReplaceAllString(string(resp), "")

	err = goluhn.Validate(orderNum)
	if err != nil {
		WriteJSON(w, http.StatusUnprocessableEntity, response.Response{
			Status:  "failed",
			Message: "Incorrect order number format",
		})
	}

	err = o.service.AddOrdersNumber(r.Context(), orderNum)
	if err != nil {
		HandleErrors(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, response.Response{
		Status:  "success",
		Message: orderNum,
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
