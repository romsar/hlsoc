package grpc

import (
	"fmt"
	"github.com/romsar/hlsoc"
	grpcgen "github.com/romsar/hlsoc/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	addr string

	s *grpc.Server

	tokenizer      hlsoc.Tokenizer
	userRepository hlsoc.UserRepository
	passwordHasher hlsoc.PasswordHasher
	postRepository hlsoc.PostRepository

	grpcgen.UnimplementedUserServiceServer
	grpcgen.UnimplementedPostServiceServer
}

type Option func(s *Server)

func WithTokenizer(tokenizer hlsoc.Tokenizer) Option {
	return func(s *Server) {
		s.tokenizer = tokenizer
	}
}

func WithUserRepository(repo hlsoc.UserRepository) Option {
	return func(s *Server) {
		s.userRepository = repo
	}
}

func WithPasswordHasher(ph hlsoc.PasswordHasher) Option {
	return func(s *Server) {
		s.passwordHasher = ph
	}
}

func WithPostRepository(repo hlsoc.PostRepository) Option {
	return func(s *Server) {
		s.postRepository = repo
	}
}

func New(addr string, opts ...Option) *Server {
	s := &Server{addr: addr}
	s.s = grpc.NewServer(grpc.ChainUnaryInterceptor(s.loggingInterceptor, s.authInterceptor))

	for _, opt := range opts {
		opt(s)
	}

	reflection.Register(s.s)

	{
		grpcgen.RegisterUserServiceServer(s.s, s)
		grpcgen.RegisterPostServiceServer(s.s, s)
	}

	return s
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("start grpc server error: %w", err)
	}

	return s.s.Serve(lis)
}

func (s *Server) Stop() {
	s.s.GracefulStop()
}
