package handler

import (
	"encoding/json"
	"errors"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"go.uber.org/zap"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Log.Info("Json encode error", zap.Error(err))
		return
	}
}

func HandleErrors(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := err.Error()

	switch {
	case errors.Is(err, apperrors.LoginAlreadyExistsErr):
		statusCode = http.StatusConflict
	case errors.Is(err, apperrors.UserNotFoundErr):
		statusCode = http.StatusUnauthorized
		message = "Incorrect login or password"
	case errors.Is(err, apperrors.OrderAddedAnotherUserErr):
		statusCode = http.StatusConflict
	case errors.Is(err, apperrors.OrderAddedThisUserErr):
		statusCode = http.StatusOK
	}

	WriteJSON(w, statusCode, response.Response{
		Message: message,
		Status:  "failed",
	})
}
