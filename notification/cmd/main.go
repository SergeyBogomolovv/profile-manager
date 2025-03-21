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
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/broker"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/repo"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/service"
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

	mailer := mailer.New(conf.SMTP)
	userRepo := repo.New(postgres)
	svc := service.New(mailer, userRepo)
	broker := broker.MustNew(logger, amqpConn, svc)

	app := app.New(logger, conf)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	broker.Consume(ctx)
	app.Start()
	<-ctx.Done()
	app.Stop()
	broker.Close()
}

func init() {
	godotenv.Load()
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
