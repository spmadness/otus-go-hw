package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func UnaryServerRequestLoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, r interface{}, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		result, err := h(ctx, r)
		if err != nil {
			return nil, err
		}

		duration := time.Since(start)
		logger.Info(fmt.Sprintf("grpc request method: %s, duration: %s", i.FullMethod, duration))

		return result, nil
	}
}
