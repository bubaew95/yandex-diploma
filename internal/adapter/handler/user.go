package handler

import (
	"encoding/json"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/authdto"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/userrequest"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/resplogindto"
	"github.com/bubaew95/yandex-diploma/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/bubaew95/yandex-diploma/internal/core/ports"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req authdto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Info("Json encode error", zap.Error(err))
		utils.WriteJSON(w, http.StatusBadRequest, response.Response{
			Status:  "failed",
			Message: "json encode error",
		})
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		logger.Log.Info("Validation error", zap.Any("errors", validationErrors))
		utils.WriteJSON(w, http.StatusBadRequest, response.Response{
			Errors: validationErrors,
			Status: "failed",
		})
		return
	}

	token, err := u.service.Registration(r.Context(), req)
	if err != nil {
		logger.Log.Info("Registration error", zap.Error(err))
		HandleErrors(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	utils.WriteJSON(w, http.StatusOK, response.Response{
		Status:  "success",
		Message: "User successfully registered and authenticated",
	})
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var signIn authdto.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&signIn); err != nil {
		logger.Log.Info("Json encode error", zap.Error(err))
		utils.WriteJSON(w, http.StatusBadRequest, response.Response{
			Status:  "failed",
			Message: "json encode error",
		})
		return
	}

	if validationErrors := signIn.Validate(); len(validationErrors) > 0 {
		logger.Log.Info("Validation error", zap.Any("errors", validationErrors))
		utils.WriteJSON(w, http.StatusBadRequest, response.Response{
			Errors: validationErrors,
			Status: "failed",
		})
		return
	}

	token, err := u.service.Authorization(r.Context(), signIn)
	if err != nil {
		logger.Log.Info("Authorization error", zap.Error(err))
		HandleErrors(w, err)
		return
	}

	tokenExires := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  tokenExires,
		HttpOnly: true,
	})

	utils.WriteJSON(w, http.StatusOK, resplogindto.ResponseToken{
		Token:  token,
		Expire: tokenExires,
	})
}

func (u UserHandler) Balance(w http.ResponseWriter, r *http.Request) {
	balance, err := u.service.Balance(r.Context())
	if err != nil {
		logger.Log.Info("Balance error", zap.Error(err))
		HandleErrors(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, balance)
}

func (u UserHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var withdraw userrequest.Withdraw
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		logger.Log.Info("Json encode error", zap.Error(err))
		utils.WriteJSON(w, http.StatusBadRequest, response.Response{
			Status:  "failed",
			Message: "json encode error",
		})
		return
	}

	err := u.service.BalanceWithdraw(r.Context(), withdraw)
	if err != nil {
		HandleErrors(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response.Response{
		Status:  "success",
		Message: "User successfully withdrawn",
	})
}

func (u UserHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	widthdraws, err := u.service.GetWithdrawals(r.Context())
	if err != nil {
		logger.Log.Info("GetWithdrawals error", zap.Error(err))
		HandleErrors(w, err)
		return
	}

	status := http.StatusOK
	if len(widthdraws) == 0 {
		status = http.StatusNoContent
	}

	utils.WriteJSON(w, status, widthdraws)
}
