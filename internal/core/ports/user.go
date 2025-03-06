package ports

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request"
	"github.com/bubaew95/yandex-diploma/internal/core/model"
)

type UserService interface {
	Registration(ctx context.Context, s request.SignUpRequest) error
}

type UserRepository interface {
	Registration(ctx context.Context, s model.SignUp) error
}
