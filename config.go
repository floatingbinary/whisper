package whisper

import (
	"context"
	"encoding/json"
)

type EventBus struct {
	// Context is the context for the event client
	Context context.Context

	// Connection is the connection string to the event bus
	Connection string

	// EventHandlers is a list of event handlers
	EventHandlers []EventHandler

	// eventPool is a map of event name to event handler
	eventPool map[Event]EventHandler
}

func NewEventBus(ctx context.Context, conn string) *EventBus {
	c := &EventBus{
		Context:    ctx,
		Connection: conn,
		eventPool:  make(map[Event]EventHandler),
	}
	return c
}

func (c *EventBus) RegisterEvents(events ...EventHandler) {
	for _, e := range events {
		c.eventPool[e.GetEventName()] = e
	}
}

func (c *EventBus) Dispatch(ctx context.Context, event Event, payload []byte) error {
	handler, ok := c.eventPool[event]
	if !ok {
		return ErrEventNotFound
	}

	var msg EventPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		return ErrInvalidPayload
	}
	if err := handler.ValidatePayload(msg.Payload.Data); err != nil {
		return ErrInvalidPayload
	}
	return handler.Handle(ctx, msg.Payload.Data)
}
