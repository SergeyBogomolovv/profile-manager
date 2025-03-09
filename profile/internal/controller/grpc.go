package controller

import (
	"context"
	"log/slog"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileService interface {
}

type gRPCController struct {
	pb.UnimplementedProfileServer
	svc      ProfileService
	logger   *slog.Logger
	validate *validator.Validate
}

func NewGRPCController(logger *slog.Logger, svc ProfileService) *gRPCController {
	validate := validator.New()
	return &gRPCController{svc: svc, logger: logger, validate: validate}
}

func (c *gRPCController) Init(srv *grpc.Server) {
	pb.RegisterProfileServer(srv, c)
}

func (c *gRPCController) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}

func (c *gRPCController) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
