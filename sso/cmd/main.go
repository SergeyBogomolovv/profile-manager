package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SergeyBogomolovv/profile-manager/common/postgres"
	"github.com/SergeyBogomolovv/profile-manager/common/redis"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/app"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/controller"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/repo"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	confPath := flag.String("config", "./config/sso.yml", "path to config file")
	flag.Parse()
	conf := config.MustLoadConfig(*confPath)

	redis := redis.MustNew(conf.RedisURL)
	postgres := postgres.MustNew(conf.PostgresURL)

	userRepo := repo.NewUserRepo(postgres)
	tokenRepo := repo.NewTokensRepo(redis)
	txManager := transaction.NewTxManager(postgres)

	authSvc := service.NewAuthService(txManager, userRepo, tokenRepo, []byte(conf.JWT.SecretKey))

	logger := newLogger()
	grpcController := controller.NewGRPCController(logger, authSvc)
	oauthController := controller.NewOAuthController(logger, conf.OAuth, authSvc)

	app := app.New(logger, conf, oauthController, grpcController)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		app.Stop()
		redis.Close()
		postgres.Close()
	}()

	app.Start()
	wg.Wait()
}

func init() {
	godotenv.Load()
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
