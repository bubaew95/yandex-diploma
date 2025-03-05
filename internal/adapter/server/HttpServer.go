package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const defaultHost = "0.0.0.0"

type HttpServer interface {
	Start()
	Stop()
}

type httpServer struct {
	Port   int
	Server *http.Server
}

func NewHttpServer(r *chi.Mux, c conf.Config) HttpServer {
	return &httpServer{
		Port: c.Port,
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", defaultHost, c.Port),
			Handler: r,
		},
	}
}

func (s *httpServer) Start() {
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("Http server start error", zap.Error(err))
		}
	}()
}

func (s *httpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Log.Info("Http server shutdown error", zap.Error(err))
	}
}
