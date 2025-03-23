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
	"github.com/SergeyBogomolovv/profile-manager/common/redis"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/app"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/broker"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/controller"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/repo"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/service"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/telegram"
	"github.com/SergeyBogomolovv/profile-manager/notification/pkg/bot"
	"github.com/joho/godotenv"
)

func main() {
	confPath := flag.String("config", "./config/profile.yml", "path to config file")
	flag.Parse()
	conf := config.MustLoadConfig(*confPath)

	redis := redis.MustNew(conf.RedisURL)
	defer redis.Close()
	postgres := postgres.MustNew(conf.PostgresURL)
	defer postgres.Close()
	amqpConn := rabbitmq.MustNew(conf.RabbitmqURL)
	defer amqpConn.Close()

	bot := bot.MustNew(conf.TelegramToken)
	logger := newLogger()

	mailer := mailer.New(conf.SMTP)
	sender := telegram.NewSender(bot)
	userRepo := repo.NewUserRepo(postgres)
	tokenRepo := repo.NewTokenRepo(redis)
	subscriptionRepo := repo.NewSubscriptionRepo(postgres)
	txManager := transaction.NewTxManager(postgres)
	svc := service.New(txManager, mailer, sender, userRepo, tokenRepo, subscriptionRepo)

	loginer := telegram.NewLoginer(logger, bot, svc)
	loginer.Init()

	broker := broker.MustNew(logger, amqpConn, svc)

	controller := controller.New(svc)
	app := app.New(logger, conf, controller, bot, broker)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	app.Start(ctx)
	<-ctx.Done()
	app.Stop()
}

func init() {
	godotenv.Load()
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
