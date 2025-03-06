package service

import (
	"context"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request/authdto"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
	"github.com/bubaew95/yandex-diploma/internal/core/ports"
	"github.com/bubaew95/yandex-diploma/internal/core/token"
	"github.com/bubaew95/yandex-diploma/pkg/crypto"
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
