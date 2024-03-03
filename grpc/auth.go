package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
)

var authMethods []string

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func (s *Server) authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if slices.Contains(authMethods, filepath.Base(info.FullMethod)) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		if !s.valid(md["authorization"]) {
			return nil, errInvalidToken
		}
	}

	return handler(ctx, req)
}

func (s *Server) valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	_, err := s.tokenizer.Verify(token)
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	return true
}
