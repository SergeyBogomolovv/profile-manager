package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const UserIdKey = "user_id"

func ExtractUserID(ctx context.Context) string {
	vals := metadata.ValueFromIncomingContext(ctx, UserIdKey)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func NewJwtInterceptor(secret []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader, exists := md["authorization"]
		if !exists || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is missing")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		if tokenString == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}

		claims, err := VerifyJWT(tokenString, secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		md.Append(UserIdKey, claims.UserID)
		ctx = metadata.NewIncomingContext(ctx, md)
		return handler(ctx, req)
	}
}
