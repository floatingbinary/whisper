package whisper

import (
	"context"
)

type RedisPubSub struct {
}

func NewRedisPubSub() *RedisPubSub {
	return &RedisPubSub{}
}

func (*RedisPubSub) Connect(ctx context.Context, conn string) error {
	return nil
}

func (c *RedisPubSub) Publish(ctx context.Context, topic string, msg []byte) error {
	return nil
}

func (c *RedisPubSub) Subscribe(dispatch EventDispatcher, events ...EventHandler) error {
	return nil
}

func (c *RedisPubSub) Close() error {
	return nil
}
