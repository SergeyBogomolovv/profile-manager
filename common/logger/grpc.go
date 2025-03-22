package logger

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func LoggerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = Inject(ctx, logger)
		return handler(ctx, req)
	}
}
