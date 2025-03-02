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

type OAuthController interface {
	Init(router *chi.Mux)
}

type GRPCController interface {
	Init(srv *grpc.Server)
}

func New(logger *slog.Logger, conf *config.Config, oauthController OAuthController, gRPCController GRPCController) *app {
	router := chi.NewRouter()
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HTTP.Port),
		Handler: router,
	}
	grpcSrv := grpc.NewServer()
	gRPCController.Init(grpcSrv)
	oauthController.Init(router)
	return &app{httpSrv: httpSrv, grpcSrv: grpcSrv, logger: logger, conf: conf}
}

func (a *app) Start() {
	go func() {
		a.logger.Info("starting http server", "port", a.conf.HTTP.Port)
		if err := a.httpSrv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to run http server %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.conf.GRPC.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		a.logger.Info("starting grpc server", "port", a.conf.GRPC.Port)
		a.grpcSrv.Serve(lis)
	}()
}

func (a *app) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	a.grpcSrv.GracefulStop()
	a.logger.Info("grpc server stopped")

	a.httpSrv.Shutdown(ctx)
	a.logger.Info("http server stopped")
}
