package grpc

import (
	"context"
	"github.com/romsar/hlsoc"
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

func (s *Server) authUnaryInterceptor(
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

		user, valid := s.valid(md["authorization"])
		if !valid {
			return nil, errInvalidToken
		}

		ctx = context.WithValue(ctx, "user", user)
	}

	return handler(ctx, req)
}

type wrappedStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func newWrappedStream(ctx context.Context, s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{ctx, s}
}

func (s *Server) authStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if slices.Contains(authMethods, filepath.Base(info.FullMethod)) {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errMissingMetadata
		}

		user, valid := s.valid(md["authorization"])
		if !valid {
			return errInvalidToken
		}

		ctx := context.WithValue(ss.Context(), "user", user)

		return handler(srv, newWrappedStream(ctx, ss))
	}

	return handler(srv, ss)
}

func (s *Server) valid(authorization []string) (*hlsoc.User, bool) {
	if len(authorization) < 1 {
		return nil, false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	user, err := s.tokenizer.Verify(token)
	if err != nil {
		slog.Error(err.Error())
		return nil, false
	}

	return user, true
}
