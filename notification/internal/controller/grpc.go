package controller

import (
	"context"
	"errors"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/notification"
	"github.com/SergeyBogomolovv/profile-manager/common/auth"
	"github.com/SergeyBogomolovv/profile-manager/common/logger"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
}

type controller struct {
	pb.UnimplementedNotificationServer
	svc      Service
	validate *validator.Validate
}

func New(svc Service) *controller {
	validate := validator.New()
	return &controller{svc: svc, validate: validate}
}

func (c *controller) Init(srv *grpc.Server) {
	pb.RegisterNotificationServer(srv, c)
}

func (c *controller) GenerateTelegramToken(ctx context.Context, req *pb.GenerateTelegramTokenRequest) (*pb.GenerateTelegramTokenResponse, error) {
	userID := auth.ExtractUserID(ctx)
	if err := c.validate.Var(userID, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	token, err := c.svc.GenerateToken(ctx, userID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		logger.Extract(ctx).Error("failed to generate token", "error", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &pb.GenerateTelegramTokenResponse{Token: token}, nil
}
