package app

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/SergeyBogomolovv/profile-manager/common/auth"
	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/config"
	"google.golang.org/grpc"
)

type app struct {
	grpcSrv *grpc.Server
	logger  *slog.Logger
	conf    *config.Config
}

type GRPCController interface {
	Init(srv *grpc.Server)
}

func New(log *slog.Logger, conf *config.Config, grpcController GRPCController) *app {
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.LoggerInterceptor(log),
			auth.JwtInterceptor([]byte(conf.JwtSecret)),
		),
	)
	grpcController.Init(grpcSrv)
	return &app{grpcSrv: grpcSrv, logger: log, conf: conf}
}

func (a *app) Start() {
	go a.startGrpcServer()
}

func (a *app) startGrpcServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.conf.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	a.logger.Info("starting grpc server", "addr", lis.Addr())
	a.grpcSrv.Serve(lis)
}

func (a *app) Stop() {
	a.grpcSrv.GracefulStop()
	a.logger.Info("grpc server stopped")
}
