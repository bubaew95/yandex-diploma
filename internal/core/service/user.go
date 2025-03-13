package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/authdto"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/userrequest"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/response/responsedto"
	"github.com/bubaew95/yandex-diploma/internal/core/entity/userentity"
	apperrors "github.com/bubaew95/yandex-diploma/internal/core/errors"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"github.com/bubaew95/yandex-diploma/pkg/crypto"
	"github.com/bubaew95/yandex-diploma/pkg/token"
	"regexp"
	"strconv"
)

type UserService struct {
	repo   ports.UserRepository
	config *conf.Config
}

func NewUserService(repo ports.UserRepository, config *conf.Config) *UserService {
	return &UserService{
		repo:   repo,
		config: config,
	}
}

func (s UserService) Registration(ctx context.Context, req authdto.SignUpRequest) (string, error) {
	newCrypto := crypto.NewCrypto(s.config.SecretKey)

	password, err := newCrypto.Encode(req.Password)
	if err != nil {
		return "", err
	}

	userModel := usermodel.UserRegistration{
		Login:    req.Login,
		Password: password,
	}

	user, err := s.repo.AddUser(ctx, userModel)
	if err != nil {
		return "", err
	}

	newJwtToken := token.NewJwtToken(s.config.SecretKey)

	jwtToken, err := newJwtToken.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s UserService) Authorization(ctx context.Context, req authdto.SignInRequest) (string, error) {
	newCrypto := crypto.NewCrypto(s.config.SecretKey)
	password, err := newCrypto.Encode(req.Password)
	if err != nil {
		return "", err
	}

	userModel := usermodel.UserLogin{
		Login:    req.Login,
		Password: password,
	}

	user, err := s.repo.FindUserByLoginAndPassword(ctx, userModel)
	if err != nil {
		return "", err
	}

	newJwtToken := token.NewJwtToken(s.config.SecretKey)

	jwtToken, err := newJwtToken.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s UserService) Balance(ctx context.Context) (usermodel.Balance, error) {
	user, ok := ctx.Value("user").(userentity.User)
	if !ok {
		return usermodel.Balance{}, apperrors.ErrUserNotFound
	}

	return s.repo.GetUserBalance(ctx, user.Id)
}

func (s UserService) BalanceWithdraw(ctx context.Context, ur userrequest.Withdraw) error {
	user, ok := ctx.Value("user").(userentity.User)
	if !ok {
		return apperrors.ErrUserNotFound
	}

	reg := regexp.MustCompile("[^0-9]+")
	num := reg.ReplaceAllString(ur.Order, "")

	err := goluhn.Validate(num)
	if err != nil {
		return apperrors.ErrInvalidOrderNumber
	}

	orderNum, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return err
	}

	userModel := usermodel.Withdraw{
		OrderNumber: orderNum,
		Amount:      ur.Sum,
		UserID:      user.Id,
	}

	return s.repo.BalanceWithdraw(ctx, userModel)
}

func (s UserService) GetWithdrawals(ctx context.Context) ([]responsedto.Withdraw, error) {
	return s.repo.GetWithdrawals(ctx)
}
