package grpc

import (
	"context"
	"errors"
	"github.com/romsar/hlsoc"
	grpcgen "github.com/romsar/hlsoc/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

var _ grpcgen.PostServiceServer = (*Server)(nil)

func init() {
	authMethods = append(authMethods, "GetFeed", "CreatePost", "StreamFeed")
}

func (s *Server) GetFeed(ctx context.Context, req *grpcgen.GetFeedRequest) (*grpcgen.GetFeedResponse, error) {
	user, ok := ctx.Value("user").(*hlsoc.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "cannot get user")
	}

	filter := &hlsoc.FeedFilter{
		UserID: user.ID,
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}

	posts, err := s.postRepository.GetFeed(ctx, filter)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "get feed posts error")
	}

	postsProto := make([]*grpcgen.Post, len(posts))
	for i := range posts {
		postsProto[i] = postToProto(posts[i])
	}

	return &grpcgen.GetFeedResponse{Posts: postsProto}, nil
}

func (s *Server) CreatePost(ctx context.Context, req *grpcgen.CreatePostRequest) (*grpcgen.CreatePostResponse, error) {
	user, ok := ctx.Value("user").(*hlsoc.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "cannot get user")
	}

	post := &hlsoc.Post{
		Text:      req.GetText(),
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
	}

	err := s.postRepository.CreatePost(ctx, post)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "create post error")
	}

	return &grpcgen.CreatePostResponse{Post: postToProto(post)}, nil
}

func (s *Server) StreamFeed(_ *grpcgen.StreamFeedRequest, stream grpcgen.PostService_StreamFeedServer) error {
	ctx := stream.Context()

	user, ok := ctx.Value("user").(*hlsoc.User)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "cannot get user")
	}

	err := s.postBroker.ConsumeNewPost(ctx, user.ID, func(post *hlsoc.Post) error {
		if err := stream.Send(postToProto(post)); err != nil {
			return err
		}
		return nil
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error(err.Error())
		return status.Errorf(codes.Internal, "internal error")
	}

	return nil
}

func postToProto(post *hlsoc.Post) *grpcgen.Post {
	return &grpcgen.Post{
		Id:        post.ID.String(),
		Text:      post.Text,
		CreatedBy: post.CreatedBy.String(),
		CreatedAt: toProtoDate(post.CreatedAt),
	}
}
