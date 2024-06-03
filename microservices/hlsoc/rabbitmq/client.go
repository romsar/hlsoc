package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
}

func New(addr string) (*Client, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
