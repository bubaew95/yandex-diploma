package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
	"github.com/bubaew95/yandex-diploma/internal/core/service"
	"github.com/bubaew95/yandex-diploma/internal/infra/repository/mock"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	type Want struct {
		StatusCode int
		Result     userentity.User
	}

	tests := []struct {
		Name string
		Data string
		Want Want
	}{
		{
			Name: "Success registration",
			Data: `{"login": "test", "password": "test"}`,
			Want: Want{
				StatusCode: http.StatusOK,
				Result: userentity.User{
					Id:    1,
					Login: "test",
				},
			},
		},
		//{
		//	Name: "Error format",
		//	Data: `{"log": "test", "password": "test"}`,
		//	Want: Want{
		//		StatusCode: http.StatusBadRequest,
		//		Result:     userentity.User{},
		//	},
		//},
		//{
		//	Name: "Login exists",
		//	Data: `{"login": "test", "password": "test"}`,
		//	Want: Want{
		//		StatusCode: http.StatusConflict,
		//		Result:     userentity.User{},
		//	},
		//},
	}

	config := conf.NewConfig()
	route := chi.NewRouter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(userRepositoryMock, config)
	userHandler := NewUserHandler(userService)

	route.Route("/api/user", func(lr chi.Router) {
		lr.Post("/register", userHandler.SignUp)
		lr.Post("/login", userHandler.Login)
	})

	ts := httptest.NewServer(route)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var registrationData usermodel.UserRegistration
			err := json.Unmarshal([]byte(tt.Data), &registrationData)
			require.NoError(t, err)

			userRepositoryMock.
				EXPECT().
				AddUser(context.Background(), registrationData).
				Return(tt.Want.Result, nil)

			req, err := http.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(tt.Data))
			require.NoError(t, err)
			defer req.Body.Close()

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			fmt.Println("code", resp.StatusCode)

			//fmt.Println("test", respBody, tt.Want.Result)
			//
			//assert.Equal(t, respBody, tt.Want.Result)
		})
	}
}
