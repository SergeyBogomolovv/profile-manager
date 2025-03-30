package main

import (
	"log"
	"log/slog"
	"net/http"

	notificationPb "github.com/SergeyBogomolovv/profile-manager/common/api/notification"
	profilePb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	ssoPb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"

	_ "github.com/SergeyBogomolovv/profile-manager/gateway/docs"
	"github.com/SergeyBogomolovv/profile-manager/gateway/internal/controller"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	ssoConn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("failed to connect to sso %v", err)
	}
	defer ssoConn.Close()
	ssoClient := ssoPb.NewSSOClient(ssoConn)

	profileConn, err := grpc.NewClient("localhost:50052", opts...)
	if err != nil {
		log.Fatalf("failed to connect to profile %v", err)
	}
	defer profileConn.Close()
	profileClient := profilePb.NewProfileClient(profileConn)

	notificationConn, err := grpc.NewClient("localhost:50053", opts...)
	if err != nil {
		log.Fatalf("failed to connect to notification %v", err)
	}
	defer notificationConn.Close()
	notificationClient := notificationPb.NewNotificationClient(notificationConn)

	profileController := controller.NewProfileController(slog.Default(), profileClient)
	authController := controller.NewAuthController(slog.Default(), ssoClient)
	notiController := controller.NewNotificationController(notificationClient)

	r := chi.NewRouter()

	r.Get("/docs/*", httpSwagger.WrapHandler)
	authController.Init(r)
	profileController.Init(r)
	notiController.Init(r)

	log.Println("starting server")
	http.ListenAndServe(":8081", r)
}

func init() {
	godotenv.Load()
}
