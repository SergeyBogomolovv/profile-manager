package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("failed to connect to addr %v", err)
	}
	defer conn.Close()

	client := pb.NewSSOClient(conn)

	resp, err := client.Login(context.Background(), &pb.LoginRequest{Email: "test@email.com", Password: "password"})
	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}

	fmt.Println(resp.AccessToken)
}
