package grpc

import (
	"context"
	"github.com/romsar/hlsoc"
	grpcgen "github.com/romsar/hlsoc/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var _ grpcgen.PostServiceServer = (*Server)(nil)

func init() {
	authMethods = append(authMethods, "GetFeed", "CreatePost")
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
		return nil, status.Errorf(codes.Internal, "create post error")
	}

	return &grpcgen.CreatePostResponse{Post: postToProto(post)}, nil
}

func postToProto(post *hlsoc.Post) *grpcgen.Post {
	return &grpcgen.Post{
		Id:        post.ID.String(),
		Text:      post.Text,
		CreatedBy: post.CreatedBy.String(),
		CreatedAt: toProtoDate(post.CreatedAt),
	}
}
