package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient("localhost:50052", opts...)
	if err != nil {
		log.Fatalf("failed to connect to addr %v", err)
	}
	defer conn.Close()

	client := pb.NewProfileClient(conn)

	body, err := os.ReadFile("image.jpeg")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.UpdateProfile(ctx, &pb.UpdateProfileRequest{Avatar: body})
	if err != nil {
		log.Fatalf("failed to update profile: %v", err)
	}

	log.Println(resp.Avatar)
}

func init() {
	godotenv.Load()
}
