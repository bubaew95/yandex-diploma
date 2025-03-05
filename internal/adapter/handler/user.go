package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/bubaew95/yandex-diploma/internal/core/ports"
)

type UserHandler struct {
	service ports.UserService
	router  *chi.Mux
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u UserHandler) InitRoute() {
	u.router.Route("/user", func(r chi.Router) {
		r.Post("/register", u.Registration)
		r.Post("/login", u.Login)
	})
}

func (u UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}
