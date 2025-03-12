package controller

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (domain.Profile, error)
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
	if err := c.validate.Var(req.UserId, "required,uuid"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	profile, err := c.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, domain.ErrProfileNotFound) {
			return nil, status.Errorf(codes.NotFound, "profile not found")
		}
		c.logger.Error("failed to get profile", "error", err)
		return nil, status.Error(codes.Internal, "failed to get profile")
	}
	return domainToGRPC(profile), nil
}

func (c *gRPCController) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}

func domainToGRPC(profile domain.Profile) *pb.ProfileResponse {
	return &pb.ProfileResponse{
		UserId:    profile.UserID,
		Username:  profile.Username,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		BirthDate: profile.BirthDate,
		Gender:    string(profile.Gender),
		Avatar:    profile.Avatar,
	}
}
