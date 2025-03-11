package apperrors

import "errors"

type ErrorResponse error

var (
	ErrLoginAlreadyExists    ErrorResponse = errors.New("login already exists")
	ErrUserNotFound          ErrorResponse = errors.New("user not found")
	ErrTokenNotFound         ErrorResponse = errors.New("token not found")
	ErrOrderNotFound         ErrorResponse = errors.New("order not found")
	ErrOrdersEmpty           ErrorResponse = errors.New("orders empty")
	ErrOrderAddedAnotherUser ErrorResponse = errors.New("order number has already been uploaded by another user")
	ErrOrderAddedThisUser    ErrorResponse = errors.New("order number has already been uploaded by this user")
	ErrInvalidOrderNumber    ErrorResponse = errors.New("incorrect order number format")
	ErrIncorrectRequest      ErrorResponse = errors.New(" Incorrect request format")
)
