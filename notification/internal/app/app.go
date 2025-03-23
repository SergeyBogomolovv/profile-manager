package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/SergeyBogomolovv/profile-manager/common/auth"
	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/config"
	"google.golang.org/grpc"
	tele "gopkg.in/telebot.v4"
)

type Controller interface {
	Init(srv *grpc.Server)
}

type Broker interface {
	Close() error
	Consume(ctx context.Context)
}

type app struct {
	logger *slog.Logger
	conf   *config.Config
	srv    *grpc.Server
	bot    *tele.Bot
	broker Broker
}

func New(log *slog.Logger, conf *config.Config, controller Controller, bot *tele.Bot, broker Broker) *app {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.LoggerInterceptor(log),
			auth.JwtInterceptor([]byte(conf.JwtSecret)),
		),
	)
	controller.Init(srv)
	return &app{logger: log, conf: conf, srv: srv, bot: bot, broker: broker}
}

func (a *app) Start(ctx context.Context) {
	go a.startServer()
	go a.startBot()
	go a.startConsumer(ctx)
}

func (a *app) startBot() {
	a.logger.Info("starting bot", "username", a.bot.Me.Username)
	a.bot.Start()
}

func (a *app) startConsumer(ctx context.Context) {
	a.logger.Info("starting consumer")
	a.broker.Consume(ctx)
}

func (a *app) startServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.conf.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	a.logger.Info("starting grpc server", "addr", lis.Addr())
	a.srv.Serve(lis)
}

func (a *app) Stop() {
	a.srv.GracefulStop()
	a.logger.Info("grpc server stopped")
	a.bot.Stop()
	a.logger.Info("bot stopped")
	a.broker.Close()
	a.logger.Info("consumer stopped")
}
