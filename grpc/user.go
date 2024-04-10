package grpc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/romsar/hlsoc"
	grpcgen "github.com/romsar/hlsoc/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

var _ grpcgen.UserServiceServer = (*Server)(nil)

func init() {
	authMethods = append(authMethods, "GetUser", "SearchUsers")
}

func (s *Server) Login(ctx context.Context, req *grpcgen.LoginRequest) (*grpcgen.LoginResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id must be uuid")
	}

	user, err := s.userRepository.GetUser(ctx, hlsoc.UserFilter{ID: id})
	if err != nil {
		if errors.Is(err, hlsoc.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "user not found or invalid password")
		}

		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if !s.passwordHasher.CheckPasswordHash(req.Password, user.Password) {
		return nil, status.Error(codes.InvalidArgument, "user not found or invalid password")
	}

	token, err := s.tokenizer.CreateToken(user, time.Hour*24*30)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &grpcgen.LoginResponse{Token: token}, nil
}

func (s *Server) Register(ctx context.Context, req *grpcgen.RegisterRequest) (*grpcgen.RegisterResponse, error) {
	password, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	user := &hlsoc.User{
		Password:   password,
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		BirthDate:  fromProtoDate(req.BirthDate),
		Gender:     hlsoc.Gender(req.Gender),
		Biography:  req.Biography,
		City:       req.City,
	}
	err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &grpcgen.RegisterResponse{UserId: user.ID.String()}, nil
}

func (s *Server) GetUser(ctx context.Context, req *grpcgen.GetUserRequest) (*grpcgen.GetUserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id must be uuid")
	}

	user, err := s.userRepository.GetUser(ctx, hlsoc.UserFilter{ID: id})
	if err != nil {
		if errors.Is(err, hlsoc.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &grpcgen.GetUserResponse{User: userToProto(user)}, nil
}

func (s *Server) SearchUsers(ctx context.Context, req *grpcgen.SearchUserRequest) (*grpcgen.SearchUserResponse, error) {
	users, err := s.userRepository.SearchUsers(ctx, hlsoc.UserFilter{
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		OrderBy:    "id",
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	usersProto := make([]*grpcgen.User, len(users))
	for i := range users {
		usersProto[i] = userToProto(users[i])
	}

	return &grpcgen.SearchUserResponse{Users: usersProto}, nil
}

func userToProto(user *hlsoc.User) *grpcgen.User {
	return &grpcgen.User{
		Id:         user.ID.String(),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		BirthDate:  toProtoDate(user.BirthDate),
		Gender:     grpcgen.Gender(user.Gender),
		Biography:  user.Biography,
		City:       user.City,
	}
}
