package whisper

import (
	"context"
	"fmt"
)

var connectionPool = make(map[string]Client)

type Client interface {
	Connect(ctx context.Context, conn string) error
	Close() error
	Publish(ctx context.Context, topic string, msg []byte) error
	Subscribe(d EventDispatcher, subs ...EventHandler) error
}

func RegisterClient(name string, client Client) {
	connectionPool[name] = client
}

func GetClient(conn string) (Client, error) {
	client, ok := connectionPool[conn]
	if !ok {
		return nil, fmt.Errorf("client %s not found", conn)
	}
	return client, nil
}

func Listen(c Client, conf *EventBus) error {
	if err := c.Connect(conf.Context, conf.Connection); err != nil {
		return fmt.Errorf("failed to connect to client: %v", err)
	}
	defer c.Close()

	RegisterClient(conf.Connection, c)
	return c.Subscribe(conf.Dispatch, conf.EventHandlers...)
}
