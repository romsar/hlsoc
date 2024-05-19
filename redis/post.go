package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/romsar/hlsoc"
	"golang.org/x/sync/errgroup"
)

func (c *Client) GetFeed(ctx context.Context, filter *hlsoc.FeedFilter) ([]*hlsoc.Post, error) {
	// если юзер использует какие-то кастомные параметры или другую пагинацию - в кеш не ходим
	shouldUseCache := filter.Limit == 20 && filter.Offset%20 == 0
	if !shouldUseCache {
		return c.postRepository.GetFeed(ctx, filter)
	}

	keyName := userFeedKeyName(filter.UserID, filter.Offset)

	val, err := c.c.Get(ctx, keyName).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("get feed for user %s from cache err: %w", filter.UserID, err)
	}

	var posts []*hlsoc.Post
	err = json.Unmarshal([]byte(val), &posts)
	if err != nil {
		return nil, fmt.Errorf("unmarshal feed for user %s from json err: %w", filter.UserID, err)
	}

	return posts, nil
}

func (c *Client) CreatePost(ctx context.Context, post *hlsoc.Post) error {
	err := c.postRepository.CreatePost(ctx, post)
	if err != nil {
		return err
	}

	// get all friends and refresh their feed
	friends, err := c.userRepository.GetFriends(ctx, post.CreatedBy)
	if err != nil {
		return fmt.Errorf("get friends for user %s: %w", post.CreatedBy, err)
	}

	errWg, ctx := errgroup.WithContext(ctx)
	for _, friendID := range friends {
		friendID := friendID

		errWg.Go(func() error {
			err := c.refreshFeedForUser(ctx, friendID)
			if err != nil {
				return fmt.Errorf("refresh feed for user: %w", err)
			}
			return nil
		})

		errWg.Go(func() error {
			// это не очень хорошо что кеширующий слой занимается отправкой событий,
			// но у меня нет времени нефакторить
			err := c.postBroker.ProduceNewPost(ctx, friendID, post)
			if err != nil {
				return fmt.Errorf("produce new post: %w", err)
			}
			return nil
		})
	}

	return errWg.Wait()
}

func (c *Client) refreshFeedForUser(ctx context.Context, userID uuid.UUID) error {
	tagName := userFeedTagName(userID)

	// invalidate old cache
	{
		keys, err := c.c.SMembers(ctx, tagName).Result()
		if err != nil {
			return fmt.Errorf("cannot get keys for tags %s: %w", tagName, err)
		}

		keys = append(keys, tagName)

		c.c.Del(ctx, keys...)
	}

	// create new cache
	{
		limit := 20
		offset := 0

		errWg, ctx := errgroup.WithContext(ctx)

		for offset+limit <= hlsoc.FeedLimit {
			o := offset

			errWg.Go(func() error {
				posts, err := c.postRepository.GetFeed(ctx, &hlsoc.FeedFilter{
					UserID: userID,
					Limit:  limit,
					Offset: o,
				})
				if err != nil {
					return fmt.Errorf("cannot get feed, user: %s, offset: %d, err: %w", userID, o, err)
				}

				bs, err := json.Marshal(posts)
				if err != nil {
					return fmt.Errorf("cannot marshal user feed into json: %w", err)
				}

				keyName := userFeedKeyName(userID, o)

				pipe := c.c.TxPipeline()
				pipe.SAdd(ctx, tagName, keyName)
				pipe.Set(ctx, keyName, string(bs), 0)

				_, err = pipe.Exec(ctx)
				if err != nil {
					return fmt.Errorf("cannot set feed cache, user: %s, offset: %d, err: %w", userID, o, err)
				}

				return nil
			})

			offset += 20
		}

		return errWg.Wait()
	}
}

func userFeedTagName(userID uuid.UUID) string {
	return fmt.Sprintf("user.%s.feed", userID)
}

func userFeedKeyName(userID uuid.UUID, offset int) string {
	return fmt.Sprintf("user.%s.feed.%d", userID, offset)
}
