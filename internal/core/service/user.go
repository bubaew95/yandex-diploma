package service

import "github.com/bubaew95/yandex-diploma/internal/core/ports"

type UserService struct {
	repo ports.OrderRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}
