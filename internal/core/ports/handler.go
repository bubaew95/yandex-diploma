package ports

import "github.com/go-chi/chi/v5"

type RouteInitialise interface {
	InitRoute() *chi.Mux
}
