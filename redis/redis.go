package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/romsar/hlsoc"
)

type Client struct {
	c              *redis.Client
	userRepository hlsoc.UserRepository
	postRepository hlsoc.PostRepository
	postBroker     hlsoc.PostBroker
}

type Option func(c *Client)

func WithUserRepository(repo hlsoc.UserRepository) Option {
	return func(c *Client) {
		c.userRepository = repo
	}
}

func WithPostRepository(repo hlsoc.PostRepository) Option {
	return func(c *Client) {
		c.postRepository = repo
	}
}

func WithPostBroker(broker hlsoc.PostBroker) Option {
	return func(c *Client) {
		c.postBroker = broker
	}
}

func New(addr string, password string, opts ...Option) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("ping err: %w", err)
	}

	c := &Client{c: client}
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) Close() error {
	return c.Close()
}
