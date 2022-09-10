package whisper

import (
	"context"
	"encoding/json"
)

type Event string

type EventPayload struct {
	Event   Event `json:"event"`
	Payload struct {
		Version   string          `json:"version"`
		Timestamp int64           `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	} `json:"payload"`
}

type EventHandler interface {
	GetEventName() Event
	GetContext() context.Context
	GetSubscriptionID() string
	Handle(ctx context.Context, body []byte) error
	ValidatePayload(payload []byte) error
}

type EventDispatcher func(ctx context.Context, event Event, payload []byte) error
