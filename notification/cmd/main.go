package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeyBogomolovv/profile-manager/common/postgres"
	"github.com/SergeyBogomolovv/profile-manager/common/rabbitmq"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/app"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	confPath := flag.String("config", "./config/profile.yml", "path to config file")
	flag.Parse()
	conf := config.MustLoadConfig(*confPath)

	postgres := postgres.MustNew(conf.PostgresURL)
	defer postgres.Close()
	amqpConn := rabbitmq.MustNew(conf.RabbitmqURL)
	defer amqpConn.Close()

	logger := newLogger()

	app := app.New(logger, conf)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	app.Start()
	<-ctx.Done()
	app.Stop()
}

func init() {
	godotenv.Load()
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
