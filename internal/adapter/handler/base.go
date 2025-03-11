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
	case errors.Is(err, apperrors.ErrLoginAlreadyExists):
		statusCode = http.StatusConflict
	case errors.Is(err, apperrors.ErrUserNotFound):
		statusCode = http.StatusUnauthorized
		message = "Incorrect login or password"
	case errors.Is(err, apperrors.ErrOrderAddedAnotherUser):
		statusCode = http.StatusConflict
	case errors.Is(err, apperrors.ErrOrderAddedThisUser):
		statusCode = http.StatusOK
	case errors.Is(err, apperrors.ErrInvalidOrderNumber):
		statusCode = http.StatusUnprocessableEntity
		message = "Incorrect order number format"
	case errors.Is(err, apperrors.ErrOrdersEmpty):
		statusCode = http.StatusNoContent
	case errors.Is(err, apperrors.ErrIncorrectRequest):
		statusCode = http.StatusBadRequest
		message = "Incorrect request format"

	}

	utils.WriteJSON(w, statusCode, response.Response{
		Message: message,
		Status:  "failed",
	})
}
