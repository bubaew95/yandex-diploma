package handler

import (
	"encoding/json"
	"github.com/bubaew95/yandex-diploma/conf"
	localMiddleware "github.com/bubaew95/yandex-diploma/internal/adapter/handler/middleware"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
	"github.com/bubaew95/yandex-diploma/internal/core/service"
	"github.com/bubaew95/yandex-diploma/internal/infra/repository/mock"
	"github.com/bubaew95/yandex-diploma/internal/utils"
	"github.com/bubaew95/yandex-diploma/pkg/crypto"
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

func setupTestServer(t *testing.T) (*conf.Config, *mock.MockUserRepository, *httptest.Server) {
	route, config, ctrl := utils.BaseTestData(t)

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(userRepositoryMock, config)
	userHandler := NewUserHandler(userService)

	route.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", userHandler.SignUp)
			r.Post("/login", userHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(localMiddleware.AuthMiddleware(config))

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", userHandler.Balance)
				r.Post("/withdraw", userHandler.Withdraw)
			})

			r.Get("/withdrawals", userHandler.Withdrawals)
		})
	})

	ts := httptest.NewServer(route)
	t.Cleanup(ts.Close)

	return config, userRepositoryMock, ts
}

func TestUserHandlerSignUp(t *testing.T) {
	t.Parallel()

	type Want struct {
		StatusCode  int
		Response    string
		ContentType string
	}

	type mockData struct {
		Data userentity.User
		Err  error
	}

	testsData := []struct {
		Name     string
		Data     string
		Method   string
		Want     Want
		MockData mockData
	}{
		{
			Name: "Success registration",
			Data: `{"login": "test", "password": "test"}`,
			Want: Want{
				StatusCode:  http.StatusOK,
				ContentType: "application/json",
				Response:    `{"status":"success","message":"User successfully registered and authenticated"}`,
			},
			MockData: mockData{
				Data: userentity.User{
					ID:    1,
					Login: "test",
				},
				Err: nil,
			},
		},
		{
			Name: "Validation error",
			Data: `{"login": " ", "password": "test"}`,
			Want: Want{
				StatusCode:  http.StatusBadRequest,
				ContentType: "application/json",
				Response:    `{"status":"failed","errors":{"login_required":"Login is required"}}`,
			},
			MockData: mockData{
				Data: userentity.User{},
				Err:  nil,
			},
		},
		{
			Name: "Login already exists",
			Data: `{"login": "test", "password": "test"}`,
			Want: Want{
				StatusCode:  http.StatusConflict,
				ContentType: "application/json",
				Response:    `{"status":"failed","message":"login already exists"}`,
			},
			MockData: mockData{
				Data: userentity.User{},
				Err:  apperrors.ErrLoginAlreadyExists,
			},
		},
	}

	config, userRepositoryMock, ts := setupTestServer(t)

	newCrypto := crypto.NewCrypto(config.SecretKey)

	for _, tt := range testsData {
		t.Run(tt.Name, func(t *testing.T) {
			var registrationData usermodel.UserRegistration
			err := json.Unmarshal([]byte(tt.Data), &registrationData)
			require.NoError(t, err)

			jwtToken := ""
			if registrationData.Login != "" && registrationData.Password != "" {
				password, err := newCrypto.Encode(registrationData.Password)
				require.NoError(t, err)

				registrationData.Password = password

				userRepositoryMock.
					EXPECT().
					AddUser(gomock.Any(), registrationData).
					Return(tt.MockData.Data, tt.MockData.Err)

				newJwtToken := token.NewJwtToken(config.SecretKey)
				jwtToken, err = newJwtToken.GenerateToken(tt.MockData.Data)
				require.NoError(t, err)
			}

			req := utils.CreateRequest(t, ts, http.MethodPost, "/api/user/register", tt.Data, jwtToken)
			resp := utils.SendUserRequest(t, req)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if jwtToken != "" {
				user, err := req.Cookie("auth_token")
				require.NoError(t, err)
				assert.Equal(t, user.Value, jwtToken)
			}

			assert.Equal(t, tt.Want.ContentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.Want.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.Want.Response, string(respBody))
		})
	}
}

func TestUserHandlerLogin(t *testing.T) {
	t.Parallel()

	type Want struct {
		StatusCode  int
		ContentType string
		Response    string
	}
	type mockData struct {
		Data userentity.User
		Err  error
	}

	testsData := []struct {
		Name     string
		Data     string
		Want     Want
		MockData mockData
	}{
		{
			Name: "Success login",
			Data: `{"login": "test", "password": "test"}`,
			Want: Want{
				StatusCode:  http.StatusOK,
				ContentType: "application/json",
				Response:    "",
			},
			MockData: mockData{
				Data: userentity.User{},
				Err:  nil,
			},
		},
		{
			Name: "Validation error",
			Data: `{"login": "", "password": "test"}`,
			Want: Want{
				StatusCode:  http.StatusBadRequest,
				ContentType: "application/json",
				Response:    `{"status":"failed","errors":{"login_required":"Login is empty"}}`,
			},
			MockData: mockData{
				Data: userentity.User{},
				Err:  nil,
			},
		},
		{
			Name: "User not found",
			Data: `{"login": "test", "password": "test"}`,
			Want: Want{
				StatusCode:  http.StatusNotFound,
				ContentType: "application/json",
				Response:    `{"status":"failed","message":"incorrect login or password"}`,
			},
			MockData: mockData{
				Data: userentity.User{},
				Err:  apperrors.ErrIncorrectUser,
			},
		},
	}

	config, userRepositoryMock, ts := setupTestServer(t)
	newCrypto := crypto.NewCrypto(config.SecretKey)

	for _, tt := range testsData {
		t.Run(tt.Name, func(t *testing.T) {
			var signIn usermodel.UserLogin
			err := json.Unmarshal([]byte(tt.Data), &signIn)
			require.NoError(t, err)

			jwtToken := ""
			if signIn.Login != "" && signIn.Password != "" {
				password, err := newCrypto.Encode(signIn.Password)
				require.NoError(t, err)

				signIn.Password = password

				userRepositoryMock.EXPECT().
					FindUserByLoginAndPassword(gomock.Any(), signIn).
					Return(tt.MockData.Data, tt.MockData.Err)

				newJwtToken := token.NewJwtToken(config.SecretKey)
				jwtToken, err = newJwtToken.GenerateToken(tt.MockData.Data)
				require.NoError(t, err)
			}

			req := utils.CreateRequest(t, ts, http.MethodPost, "/api/user/login", tt.Data, jwtToken)
			resp := utils.SendUserRequest(t, req)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if jwtToken != "" {
				user, err := req.Cookie("auth_token")
				require.NoError(t, err)
				assert.Equal(t, user.Value, jwtToken)
			}

			if tt.Want.Response != "" {
				assert.Equal(t, tt.Want.Response, string(respBody))
			}

			assert.Equal(t, tt.Want.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.Want.ContentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestBalanceHandler(t *testing.T) {
	t.Parallel()

	type want struct {
		StatusCode  int
		ContentType string
		Result      usermodel.Balance
	}

	tests := []struct {
		Name string
		Want want
	}{
		{
			Name: "Simple response balance",
			Want: want{
				StatusCode:  http.StatusOK,
				ContentType: "application/json",
				Result: usermodel.Balance{
					Current:   100.5,
					Withdrawn: 100.5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config, userRepositoryMock, ts := setupTestServer(t)

			userRepositoryMock.EXPECT().
				GetUserBalance(gomock.Any(), gomock.Any()).
				Return(tt.Want.Result, nil)

			jwtToken := token.NewJwtToken(config.SecretKey)
			newJwtToken, err := jwtToken.GenerateToken(userentity.User{
				ID:    1,
				Login: "test",
			})
			require.NoError(t, err)

			req := utils.CreateRequest(t, ts, http.MethodGet, "/api/user/balance", "", newJwtToken)
			resp := utils.SendUserRequest(t, req)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.Want.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.Want.ContentType, resp.Header.Get("Content-Type"))

			result, err := json.Marshal(tt.Want.Result)
			require.NoError(t, err)

			assert.JSONEq(t, string(result), string(respBody))
		})
	}
}

func TestWithdrawHandler(t *testing.T) {
	t.Parallel()

	type want struct {
		StatusCode  int
		ContentType string
		Result      error
	}

	tests := []struct {
		Name        string
		Want        want
		RequestData string
		MockData    usermodel.Withdraw
	}{
		{
			Name: "Simple response withdraw",
			Want: want{
				StatusCode:  http.StatusOK,
				ContentType: "application/json",
				Result:      nil,
			},
			RequestData: `{"order": "52383003218", "sum": 100.5}`,
			MockData: usermodel.Withdraw{
				UserID:      1,
				OrderNumber: 52383003218,
				Amount:      100.5,
			},
		},
		{
			Name: "Bad request",
			Want: want{
				StatusCode:  http.StatusBadRequest,
				ContentType: "application/json",
				Result:      nil,
			},
			RequestData: "",
			MockData:    usermodel.Withdraw{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config, userRepositoryMock, ts := setupTestServer(t)

			userRepositoryMock.EXPECT().
				BalanceWithdraw(gomock.Any(), tt.MockData).
				Return(tt.Want.Result)
			jwtToken := token.NewJwtToken(config.SecretKey)
			newJwtToken, err := jwtToken.GenerateToken(userentity.User{
				ID:    1,
				Login: "test",
			})
			require.NoError(t, err)

			req := utils.CreateRequest(t, ts, http.MethodPost, "/api/user/balance/withdraw", tt.RequestData, newJwtToken)
			resp := utils.SendUserRequest(t, req)
			defer resp.Body.Close()

			assert.Equal(t, tt.Want.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.Want.ContentType, resp.Header.Get("Content-Type"))
			//assert.JSONEq(t, string(respBody), string(withdraw))
		})
	}
}
