package hlsoc

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedBy uuid.UUID `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FeedFilter struct {
	UserID        uuid.UUID
	Limit, Offset int
}

const FeedLimit = 1000

type PostRepository interface {
	GetFeed(ctx context.Context, filter *FeedFilter) ([]*Post, error)
	CreatePost(ctx context.Context, post *Post) error
}

type PostBroker interface {
	ProduceNewPost(ctx context.Context, userID uuid.UUID, post *Post) error
	ConsumeNewPost(ctx context.Context, userID uuid.UUID, f func(post *Post) error) error
}
