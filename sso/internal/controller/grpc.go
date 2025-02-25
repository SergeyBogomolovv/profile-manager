package controller

import (
	pb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"
	"google.golang.org/grpc"
)

type AuthService interface{}

type gRPCController struct {
	pb.UnimplementedSSOServer
	svc AuthService
}

func NewGRPCController(svc AuthService) *gRPCController {
	return &gRPCController{svc: svc}
}

func (c *gRPCController) Init(srv *grpc.Server) {
	pb.RegisterSSOServer(srv, c)
}
