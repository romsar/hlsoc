package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

func (s *Server) loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	rsp, err := handler(ctx, req)

	slog.DebugContext(ctx,
		"grpc: method=%s\tduration=%s\terror=%v\treq=%v\trsp%v\n",
		info.FullMethod,
		time.Since(start),
		err,
		req,
		rsp,
	)

	return rsp, err
}
