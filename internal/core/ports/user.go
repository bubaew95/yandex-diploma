package ports

import (
	"context"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/authdto"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/userrequest"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
)

type UserService interface {
	Registration(ctx context.Context, s authdto.SignUpRequest) (string, error)
	Authorization(ctx context.Context, s authdto.SignInRequest) (string, error)
	Balance(ctx context.Context) (usermodel.Balance, error)
	BalanceWithdraw(ctx context.Context, ur userrequest.Withdraw) error
}

type UserRepository interface {
	AddUser(ctx context.Context, s usermodel.UserRegistration) (userentity.User, error)
	FindUserByLoginAndPassword(ctx context.Context, s usermodel.UserLogin) (userentity.User, error)
	GetUserBalance(ctx context.Context, userID int64) (usermodel.Balance, error)
	BalanceWithdraw(ctx context.Context, ur usermodel.Withdraw) error
}
