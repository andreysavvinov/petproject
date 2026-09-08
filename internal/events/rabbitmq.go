package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &Publisher{ch: ch}, nil
}

func (p *Publisher) PublishUserDeleted(ctx context.Context, userID string, myMusicID int64) error {
	body, err := json.Marshal(UserDeleted{UserID: userID, MyMusicID: myMusicID})
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(ctx, ExchangeName, RoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
}

type Handler func(ctx context.Context, event UserDeleted) error

type Consumer struct {
	url     string
	handler Handler
}

func NewConsumer(url string, handler Handler) *Consumer {
	return &Consumer{url: url, handler: handler}
}

func (c *Consumer) Run(ctx context.Context) error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("rabbitmq dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	if err := ch.QueueBind(q.Name, RoutingKey, ExchangeName, false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	log.Printf("listening for %s events", RoutingKey)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			var event UserDeleted
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("invalid event: %v", err)
				continue
			}
			if err := c.handler(ctx, event); err != nil {
				log.Printf("handle event: %v", err)
			}
		}
	}
}
