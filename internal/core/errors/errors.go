package apperrors

import "errors"

type ErrorResponse error

var (
	LoginAlreadyExistsErr    ErrorResponse = errors.New("login already exists")
	UserNotFoundErr          ErrorResponse = errors.New("user not found")
	TokenNotFoundErr         ErrorResponse = errors.New("token not found")
	OrderNotFoundErr         ErrorResponse = errors.New("order not found")
	OrderAddedAnotherUserErr ErrorResponse = errors.New("order number has already been uploaded by another user")
	OrderAddedThisUserErr    ErrorResponse = errors.New("order number has already been uploaded by this user")
)
