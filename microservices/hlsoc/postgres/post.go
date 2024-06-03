package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/romsar/hlsoc"
)

func (db *DB) GetFeed(ctx context.Context, filter *hlsoc.FeedFilter) ([]*hlsoc.Post, error) {
	args := pgx.NamedArgs{
		"userId": filter.UserID.String(),
	}

	limit := min(filter.Limit, max(hlsoc.FeedLimit-filter.Offset-filter.Limit, 0)) // feed is limited (1000 posts max)

	query := `
		SELECT id, text, created_by, created_at 
		FROM posts
		WHERE created_by IN (
			SELECT friend_id FROM user_friends WHERE user_id = @userId
		)
		` + FormatOrderBy("created_at", "DESC") + `
		` + FormatLimitOffset(limit, filter.Offset)

	rows, err := db.db.QueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query feed: %w", err)
	}
	defer rows.Close()

	var posts []*hlsoc.Post
	for rows.Next() {
		post := hlsoc.Post{}

		err = rows.Scan(
			&post.ID,
			&post.Text,
			&post.CreatedBy,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (db *DB) CreatePost(ctx context.Context, post *hlsoc.Post) error {
	query := `
		INSERT INTO posts (id, text, created_by, created_at) 
		VALUES (@id, @text, @createdBy, @createdAt)
		RETURNING id
	`
	id := uuid.New()
	args := pgx.NamedArgs{
		"id":        id,
		"text":      post.Text,
		"createdBy": post.CreatedBy,
		"createdAt": post.CreatedAt,
	}
	_, err := db.db.ExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert post row: %w", err)
	}

	post.ID = id

	return nil
}
