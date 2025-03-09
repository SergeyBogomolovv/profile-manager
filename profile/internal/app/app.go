package app

import (
	"fmt"
	"log"
	"log/slog"
	"net"

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

func New(logger *slog.Logger, conf *config.Config, grpcController GRPCController) *app {
	grpcSrv := grpc.NewServer()
	grpcController.Init(grpcSrv)
	return &app{grpcSrv: grpcSrv, logger: logger, conf: conf}
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
