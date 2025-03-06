package apperrors

import "errors"

var (
	LoginAlreadyExists = errors.New("login already exists")
	UserNotFound       = errors.New("user not found")
)
