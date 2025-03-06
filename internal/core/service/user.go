package service

import (
	"context"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/core/dto/request"
	"github.com/bubaew95/yandex-diploma/internal/core/model"
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

func (s *UserService) Registration(ctx context.Context, req request.SignUpRequest) (string, error) {
	newCrypto := crypto.NewCrypto(s.config.SecretKey)

	password, err := newCrypto.Encode(req.Password)
	if err != nil {
		return "", err
	}

	userModel := model.SignUp{
		Login:    req.Login,
		Password: password,
	}

	user, err := s.repo.Registration(ctx, userModel)
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
