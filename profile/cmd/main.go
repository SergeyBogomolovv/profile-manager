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
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/app"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/broker"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/controller"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/repo"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/service"
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

	imageRepo := repo.MustNewImageRepo(conf.S3)
	profileRepo := repo.NewProfileRepo(postgres)
	profileSvc := service.NewProfileService(profileRepo, imageRepo)
	grpcController := controller.NewGRPCController(logger, profileSvc)

	broker := broker.MustNew(logger, amqpConn, profileSvc)

	app := app.New(logger, conf, grpcController)

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
