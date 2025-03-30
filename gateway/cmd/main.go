package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	notificationPb "github.com/SergeyBogomolovv/profile-manager/common/api/notification"
	profilePb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	ssoPb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"

	_ "github.com/SergeyBogomolovv/profile-manager/gateway/docs"
	"github.com/SergeyBogomolovv/profile-manager/gateway/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/gateway/internal/controller"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	confPath := flag.String("config", "./config/gateway.yml", "path to config file")
	flag.Parse()
	conf := config.MustLoadConfig(*confPath)

	logger := newLogger()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	ssoConn, err := grpc.NewClient(conf.SsoAddr, opts...)
	exitOnErr("failed to connect to sso", err)
	defer ssoConn.Close()
	ssoClient := ssoPb.NewSSOClient(ssoConn)

	profileConn, err := grpc.NewClient(conf.ProfileAddr, opts...)
	exitOnErr("failed to connect to profile", err)
	defer profileConn.Close()
	profileClient := profilePb.NewProfileClient(profileConn)

	notificationConn, err := grpc.NewClient(conf.NotificationAddr, opts...)
	exitOnErr("failed to connect to notification", err)
	defer notificationConn.Close()
	notificationClient := notificationPb.NewNotificationClient(notificationConn)

	profileController := controller.NewProfileController(logger, profileClient)
	authController := controller.NewAuthController(logger, ssoClient)
	notiController := controller.NewNotificationController(logger, notificationClient)

	r := chi.NewRouter()

	r.Get("/docs/*", httpSwagger.WrapHandler)
	authController.Init(r)
	profileController.Init(r)
	notiController.Init(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HttpPort),
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go srv.ListenAndServe()
	logger.Info("http server started", "addr", srv.Addr)
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	logger.Info("http server stopped")
}

func init() {
	godotenv.Load()
}

func exitOnErr(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
