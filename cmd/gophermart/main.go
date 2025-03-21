package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bubaew95/yandex-diploma/conf"
	"github.com/bubaew95/yandex-diploma/internal/adapter/handler"
	localMiddleware "github.com/bubaew95/yandex-diploma/internal/adapter/handler/middleware"
	"github.com/bubaew95/yandex-diploma/internal/adapter/logger"
	"github.com/bubaew95/yandex-diploma/internal/adapter/server"
	"github.com/bubaew95/yandex-diploma/internal/core/model/usermodel"
	"github.com/bubaew95/yandex-diploma/internal/core/service"
	"github.com/bubaew95/yandex-diploma/internal/infra"
	"github.com/bubaew95/yandex-diploma/internal/infra/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load("./.env", "./.env.local"); err != nil {
		fmt.Println("No .env file found")
	}
}

func initMigrations(db *sql.DB) error {
	return goose.Up(db, "./migrations")
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

	err = initMigrations(DB.DB)
	if err != nil {
		log.Fatalf("Initializing database migrations: %v", err)
	}

	route := chi.NewRouter()
	route.Use(localMiddleware.LoggerMiddleware)
	route.Use(middleware.AllowContentEncoding("gzip"))

	userRepository := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepository, config)
	userHandler := handler.NewUserHandler(userService)

	orderRepository := repository.NewOrdersRepository(DB)
	orderService := service.NewOrdersService(orderRepository, config)
	orderHandler := handler.NewOrdersHandler(orderService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan error, 1)
	orderService.Worker(ctx, resultCh)

	go func() {
		for res := range resultCh {
			fmt.Println("error", res)
		}
	}()

	route.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", userHandler.SignUp)
			r.Post("/login", userHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(localMiddleware.AuthMiddleware(config))
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", orderHandler.CreateOrder)
				r.Get("/", orderHandler.UserOrders)
			})

			r.Route("/balance", func(r chi.Router) {
				r.Get("/", userHandler.Balance)
				r.Post("/withdraw", userHandler.Withdraw)
			})

			r.Get("/withdrawals", userHandler.Withdrawals)
		})
	})

	runServer(route, config)

	res, err := DB.Query("SELECT balance FROM  user_balance")
	if err != nil {
		log.Fatal(err)
	}

	var users []usermodel.Withdraw
	for res.Next() {
		var user usermodel.Withdraw

		if err := res.Scan(&user.Amount); err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	if res.Err() != nil {
		log.Fatal(err)
	}

	logger.Log.Info("User balance ", zap.Any("users", users))
}

func runServer(route *chi.Mux, config *conf.Config) {
	apiRoute := chi.NewRouter()
	apiRoute.Mount("/api", route)

	httpServer := server.NewHTTPServer(apiRoute, *config)
	httpServer.Start()
	defer httpServer.Stop()

	logger.Log.Info("Server started on address", zap.String("address", config.RunAddress))
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	logger.Log.Info("Shutting down...")
}
