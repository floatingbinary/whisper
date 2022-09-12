package whisper

import (
	"context"
	"encoding/json"
)

// Event is a string that represents an event
type Event string

// EventPayload is the payload that is sent to the event handler
type EventPayload struct {
	Event   Event `json:"event"`
	Payload struct {
		Version   string          `json:"version"`
		Timestamp int64           `json:"timestamp"`
		Data      json.RawMessage `json:"data"`
	} `json:"payload"`
}

type EventHandler interface {
	// GetEventName returns the name of the event
	GetEventName() Event

	// GetSubscriptionID returns the subscription id
	GetSubscriptionID() string

	// GetContext returns the context
	GetContext() context.Context

	// ValidatePayload validates the event payload
	ValidatePayload(payload []byte) error

	// Handle handles the event
	Handle(ctx context.Context, body []byte) error
}

// EventDispatcher is a function that dispatches an event
type EventDispatcher func(ctx context.Context, event Event, payload []byte) error
