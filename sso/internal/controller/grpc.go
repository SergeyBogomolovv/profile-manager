package controller

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password, ip string) (domain.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, refreshToken string) error
}

type gRPCController struct {
	pb.UnimplementedSSOServer
	svc      AuthService
	logger   *slog.Logger
	validate *validator.Validate
}

func NewGRPCController(logger *slog.Logger, svc AuthService) *gRPCController {
	validate := validator.New()
	return &gRPCController{svc: svc, logger: logger, validate: validate}
}

func (c *gRPCController) Init(srv *grpc.Server) {
	pb.RegisterSSOServer(srv, c)
}

func (c *gRPCController) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	const op = "grpc.Register"
	logger := c.logger.With(slog.String("op", op), slog.String("email", req.Email))

	if err := c.validate.Var(req.Email, "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}
	userID, err := c.svc.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", req.Email)
		}
		logger.Error("failed to register user", "error", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &pb.RegisterResponse{UserId: userID}, nil
}

func (c *gRPCController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.TokensResponse, error) {
	const op = "grpc.Login"
	logger := c.logger.With(slog.String("op", op), slog.String("email", req.Email))
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to get peer from context")
	}

	if err := c.validate.Var(req.Email, "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	tokens, err := c.svc.Login(ctx, req.Email, req.Password, p.Addr.String())
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		logger.Error("failed to login user", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}
	return &pb.TokensResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (c *gRPCController) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AccessTokenResponse, error) {
	const op = "grpc.Refresh"
	logger := c.logger.With(slog.String("op", op))

	token, err := c.svc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		logger.Error("failed to refresh token", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to refresh token: %v", err)
	}
	return &pb.AccessTokenResponse{AccessToken: token}, nil
}

func (c *gRPCController) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	const op = "grpc.Logout"
	logger := c.logger.With(slog.String("op", op))

	if err := c.svc.Logout(ctx, req.RefreshToken); err != nil {
		logger.Error("failed to logout", "error", err)
		return nil, status.Error(codes.Internal, "failed to logout")
	}
	return &pb.LogoutResponse{Status: "success"}, nil
}
