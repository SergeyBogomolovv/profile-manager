package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/config"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type app struct {
	logger  *slog.Logger
	grpcSrv *grpc.Server
	httpSrv *http.Server
	conf    *config.Config
}

type HTTPController interface {
	Init(router *chi.Mux)
}

type GRPCController interface {
	Init(srv *grpc.Server)
}

func New(log *slog.Logger, conf *config.Config, httpController HTTPController, gRPCController GRPCController) *app {
	router := chi.NewRouter()

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HttpPort),
		Handler: router,
	}
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(logger.LoggerInterceptor(log)))

	gRPCController.Init(grpcSrv)
	httpController.Init(router)

	return &app{httpSrv: httpSrv, grpcSrv: grpcSrv, logger: log, conf: conf}
}

func (a *app) Start() {
	go a.startHTTPServer()
	go a.startGRPCServer()
}

func (a *app) Stop() {
	a.stopGRPCServer()
	a.stopHTTPServer()
}

func (a *app) startHTTPServer() {
	a.logger.Info("starting http server", "addr", a.httpSrv.Addr)
	if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to run http server %v", err)
	}
}

func (a *app) stopHTTPServer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := a.httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown http server: %v", err)
	}
	a.logger.Info("http server stopped")
}

func (a *app) startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.conf.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	a.logger.Info("starting grpc server", "addr", lis.Addr())
	a.grpcSrv.Serve(lis)
}

func (a *app) stopGRPCServer() {
	a.grpcSrv.GracefulStop()
	a.logger.Info("grpc server stopped")
}
