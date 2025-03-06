package handler

import (
	"encoding/json"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/bubaew95/yandex-diploma/internal/core/ports"
)

type UserHandler struct {
	service ports.UserService
	router  *chi.Mux
}

func NewUserHandler(r *chi.Mux, service ports.UserService) *UserHandler {
	return &UserHandler{
		service: service,
		router:  r,
	}
}

func (u UserHandler) InitRoute() {
	u.router.Route("/user", func(r chi.Router) {
		r.Post("/register", u.SignUp)
		r.Post("/login", u.Login)
	})
}

func (u UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req request.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Info("Json encode error", zap.Error(err))
		WriteJSON(w, http.StatusBadRequest, response.ErrorResponse{
			Status:  "failed",
			Message: "json encode error",
		})
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		logger.Log.Info("Validation error", zap.Any("errors", validationErrors))
		WriteJSON(w, http.StatusBadRequest, response.ErrorResponse{
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

	WriteJSON(w, http.StatusOK, response.SuccessResponseSignUp{
		Status:  "success",
		Message: "User successfully registered and authenticated",
	})
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}
