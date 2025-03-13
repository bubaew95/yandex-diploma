package handler

import (
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/adapter/handler/middleware"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/ordersmodel"
	"github.com/bubaew95/yandex-diploma/internal/core/service"
	"github.com/bubaew95/yandex-diploma/internal/infra/repository/mock"
	"github.com/bubaew95/yandex-diploma/internal/utils"
	"github.com/bubaew95/yandex-diploma/pkg/token"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupOrderTestServer(t *testing.T) (*conf.Config, *mock.MockOrderRepository, *httptest.Server) {
	route, config, ctrl := utils.BaseTestData(t)

	orderRepositoryMock := mock.NewMockOrderRepository(ctrl)
	orderService := service.NewOrdersService(orderRepositoryMock, config)
	orderHandler := NewOrdersHandler(orderService)

	route.Use(middleware.AuthMiddleware(config))
	route.Route("/api/user/orders", func(lr chi.Router) {
		lr.Post("/", orderHandler.CreateOrder)
		lr.Get("/", orderHandler.UserOrders)
	})

	ts := httptest.NewServer(route)
	t.Cleanup(ts.Close)

	return config, orderRepositoryMock, ts
}

func TestOrdersHandlerCreateOrder(t *testing.T) {
	t.Parallel()

	type want struct {
		StatusCode  int
		Result      string
		ContentType string
	}

	type mockData struct {
		Data ordersmodel.Order
		Err  error
	}

	testsData := []struct {
		Name     string
		Data     string
		Want     want
		MockData mockData
	}{
		{
			Name: "Order number add success",
			Data: `5062821234567892`,
			Want: want{
				StatusCode:  http.StatusAccepted,
				Result:      `{"status":"success", "message":"5062821234567892"}`,
				ContentType: "application/json",
			},
			MockData: mockData{
				Data: ordersmodel.Order{
					Number: 5062821234567892,
					UserId: 1,
				},
				Err: nil,
			},
		},
		{
			Name: "Invalid order number",
			Data: `235235`,
			Want: want{
				StatusCode:  http.StatusUnprocessableEntity,
				Result:      `{"status":"failed","message":"Incorrect order number format"}`,
				ContentType: "application/json",
			},
			MockData: mockData{
				Data: ordersmodel.Order{},
				Err:  nil,
			},
		},
		{
			Name: "Already been uploaded by another user",
			Data: `5062821234567892`,
			Want: want{
				StatusCode:  http.StatusConflict,
				Result:      `{"status":"failed","message":"order number has already been uploaded by another user"}`,
				ContentType: "application/json",
			},
			MockData: mockData{
				Data: ordersmodel.Order{
					Number: 5062821234567892,
					UserId: 1,
				},
				Err: apperrors.ErrOrderAddedAnotherUser,
			},
		},
		{
			Name: "Already been uploaded by this user",
			Data: `5062821234567892`,
			Want: want{
				StatusCode:  http.StatusOK,
				Result:      `{"status":"failed", "message": "order number has already been uploaded by this user"}`,
				ContentType: "application/json",
			},
			MockData: mockData{
				Data: ordersmodel.Order{
					Number: 5062821234567892,
					UserId: 1,
				},
				Err: apperrors.ErrOrderAddedThisUser,
			},
		},
		{
			Name: "Empty order number",
			Data: ``,
			Want: want{
				StatusCode:  http.StatusBadRequest,
				Result:      `{"status":"failed", "message": "Incorrect request format"}`,
				ContentType: "application/json",
			},
			MockData: mockData{
				Data: ordersmodel.Order{},
				Err:  nil,
			},
		},
	}
	config, orderRepositoryMock, ts := setupOrderTestServer(t)

	for _, tt := range testsData {
		t.Run(tt.Name, func(t *testing.T) {

			orderRepositoryMock.EXPECT().
				AddOrdersNumber(gomock.Any(), tt.MockData.Data).
				Return(tt.MockData.Err)

			jwtToken := token.NewJwtToken(config.SecretKey)
			newJwtToken, err := jwtToken.GenerateToken(userentity.User{
				Id:    1,
				Login: "test",
			})
			require.NoError(t, err)

			req := utils.CreateRequest(t, ts, http.MethodPost, "/api/user/orders", tt.Data, newJwtToken)
			resp := utils.SendUserRequest(t, req)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.JSONEq(t, tt.Want.Result, string(respBody))
			assert.Equal(t, tt.Want.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.Want.ContentType, resp.Header.Get("Content-Type"))
		})
	}
}
