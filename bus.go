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
	// only register events that are not already registered
	for _, e := range events {
		if _, ok := c.eventPool[e.GetEventName()]; !ok {
			c.eventPool[e.GetEventName()] = e
			c.EventHandlers = append(c.EventHandlers, e)
		}
	}
}

func (c *EventBus) Dispatch(ctx context.Context, e Event, payload []byte) error {
	handler, ok := c.eventPool[e]
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
