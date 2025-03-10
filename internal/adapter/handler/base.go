package handler

import (
	"errors"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/utils"
	"net/http"
)

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
	case errors.Is(err, apperrors.InvalidOrderNumberErr):
		statusCode = http.StatusUnprocessableEntity
		message = "Incorrect order number format"
	case errors.Is(err, apperrors.OrdersEmptyErr):
		statusCode = http.StatusNoContent
	case errors.Is(err, apperrors.IncorrectRequestErr):
		statusCode = http.StatusBadRequest
		message = "Incorrect request format"

	}

	utils.WriteJSON(w, statusCode, response.Response{
		Message: message,
		Status:  "failed",
	})
}
