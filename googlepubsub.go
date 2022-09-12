package whisper

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
)

type GooglePubSub struct {
	*pubsub.Client
}

func NewGooglePubSub() *GooglePubSub {
	return &GooglePubSub{}
}

func (g *GooglePubSub) Connect(ctx context.Context, conn string) error {
	c, err := pubsub.NewClient(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	g.Client = c
	return nil
}

func (c *GooglePubSub) Publish(ctx context.Context, topic string, msg []byte) error {
	t := c.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (c *GooglePubSub) Subscribe(dispatch EventDispatcher, events ...EventHandler) error {
	defer c.Close()

	var errMessages []string
	for _, e := range events {
		sub := c.Subscription(e.GetSubscriptionID())
		err := sub.Receive(e.GetContext(), func(ctx context.Context, m *pubsub.Message) {
			err := dispatch(ctx, e.GetEventName(), m.Data)
			if err == ErrInvalidPayload {
				fmt.Printf("invalid payload: %v\n", err)
				m.Ack()
			} else if err == ErrEventNotFound {
				fmt.Printf("event not found: %v\n", err)
				m.Nack()
			} else if err != nil {
				fmt.Printf("failed to handle event: %v\n", err)
				m.Nack()
			} else {
				m.Ack()
			}
		})
		if err != nil {
			fmt.Printf("failed to receive message: %v\n", err)
			errMessages = append(errMessages, fmt.Sprintf("failed to receive message: %v", err))
		}
	}
	if len(errMessages) > 0 {
		return fmt.Errorf("failed to subscribe: %v", strings.Join(errMessages, "\n"))
	}
	return nil
}
