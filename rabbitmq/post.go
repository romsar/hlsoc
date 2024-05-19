package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romsar/hlsoc"
	"time"
)

func (c *Client) ProduceNewPost(ctx context.Context, userID uuid.UUID, post *hlsoc.Post) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}

	bs, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("marshal json post: %w", err)
	}

	err = ch.ExchangeDeclare(
		"new-post",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("decrale exchange: %w", err)
	}

	err = ch.PublishWithContext(ctx,
		"new-post",
		"user."+userID.String(),
		false,
		false,
		amqp.Publishing{
			Body:       bs,
			Expiration: fmt.Sprintf("%d", (10 * time.Minute).Milliseconds()),
		},
	)
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}

func (c *Client) ConsumeNewPost(ctx context.Context, userID uuid.UUID, f func(post *hlsoc.Post) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		"new-post",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("decrale exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("decrale queue: %w", err)
	}

	err = ch.QueueBind(
		q.Name,
		"user."+userID.String(),
		"new-post",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}

	msgCh, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-msgCh:
			var post *hlsoc.Post
			err = json.Unmarshal(msg.Body, &post)
			if err != nil {
				return fmt.Errorf("unmarshal json post: %w", err)
			}

			err = f(post)
			if err != nil {
				return err
			}
		}
	}
}
