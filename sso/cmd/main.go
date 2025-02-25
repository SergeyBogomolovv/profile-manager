package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SergeyBogomolovv/profile-manager/sso/internal/app"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/controller"
	"github.com/joho/godotenv"
)

func main() {
	confPath := flag.String("config", "./config/sso.yml", "path to config file")
	flag.Parse()
	conf := config.MustLoadConfig(*confPath)

	logger := newLogger()
	grpcController := controller.NewGRPCController(nil)
	oauthController := controller.NewOAuthController(conf.OAuth, nil)

	app := app.New(logger, conf, oauthController, grpcController)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		app.Stop()
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
