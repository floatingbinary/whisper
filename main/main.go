package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/10hourlabs/whisper"
)

type HelloWorldPayload struct {
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Message   string `json:"message"`
}

type HelloWorldEvent struct{}

func (*HelloWorldEvent) GetEventName() whisper.Event {
	return "hello-world"
}

func (*HelloWorldEvent) GetSubscriptionID() string {
	return "hello-world-sub-id"
}

func (*HelloWorldEvent) GetContext() context.Context {
	return context.Background()
}

func (*HelloWorldEvent) ValidatePayload(payload []byte) error {
	var p HelloWorldPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return whisper.ErrInvalidPayload
	}
	return nil
}

func (*HelloWorldEvent) Handle(ctx context.Context, body []byte) error {
	var data HelloWorldPayload
	json.Unmarshal(body, &data) // gauranteed to not error
	fmt.Printf("%s\r ", data.Message)
	return nil
}

func main() {
	bus := whisper.NewEventBus(context.Background(), "connection-string")
	bus.RegisterEvents(&HelloWorldEvent{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := whisper.Listen(bus, whisper.NewGooglePubSub()); err != nil {
			log.Fatalf("failed to subscribe: %v", err)
		}
		wg.Done()
	}()
	wg.Wait()
}
