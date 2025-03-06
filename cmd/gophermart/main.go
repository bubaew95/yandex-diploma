package main

import (
	"fmt"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/adapter/handler"
	localMiddleware "github.com/bubaew95/yandex-diploma/internal/adapter/handler/middleware"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/bubaew95/yandex-diploma/internal/adapter/server"
	"github.com/bubaew95/yandex-diploma/internal/core/service"
	"github.com/bubaew95/yandex-diploma/internal/infra"
	"github.com/bubaew95/yandex-diploma/internal/infra/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	config := conf.NewConfig()
	if err := logger.Load(); err != nil {
		log.Fatalf("loading config: %v", err)
	}

	DB, err := infra.NewDB(config)
	if err != nil {
		log.Fatalf("Opening database connection: %v", err)
	}

	route := chi.NewRouter()
	route.Use(localMiddleware.LoggerMiddleware)
	route.Use(middleware.AllowContentEncoding("gzip"))

	userRepository := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepository, config)
	useHandler := handler.NewUserHandler(route, userService)
	useHandler.InitRoute()

	runServer(route, config)
}

func runServer(route *chi.Mux, config *conf.Config) {
	apiRoute := chi.NewRouter()
	apiRoute.Mount("/api", route)

	httpServer := server.NewHttpServer(apiRoute, *config)
	httpServer.Start()
	defer httpServer.Stop()

	logger.Log.Info("Server started on address", zap.String("address", config.RunAddress))
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	logger.Log.Info("Shutting down...")
}
