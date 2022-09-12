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
	Subscribe(dispatch EventDispatcher, handlers ...EventHandler) error
}

func RegisterClient(name string, c Client) {
	connectionPool[name] = c
}

func GetClient(conn string) (Client, error) {
	client, ok := connectionPool[conn]
	if !ok {
		return nil, fmt.Errorf("client %s not found", conn)
	}
	return client, nil
}

func Listen(e *EventBus, c Client) error {
	if err := c.Connect(e.Context, e.Connection); err != nil {
		return fmt.Errorf("failed to connect to client: %v", err)
	}
	defer c.Close()

	RegisterClient(e.Connection, c)
	return c.Subscribe(e.Dispatch, e.EventHandlers...)
}
