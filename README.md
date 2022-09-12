# whisper

Simple implementation of an Event Bus using the Publish/Subscribe pattern. Whisper comes with default implementations for [Google Pub/Sub](https://cloud.google.com/pubsub/docs/overview) and [Redis Pub/Sub](https://redis.io/docs/manual/pubsub/). It also provides a simple interface for implementing your own Pub/Sub for example using [Kafka](https://kafka.apache.org/documentation/).

## Qucik Overview

Whisper's event bus allows [publish/subscribe-style](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern) communication between your microservices without requiring the components to explicitly be aware of each other, as shown in the following diagram:

![event-pubsub](./event-driven-communication.png)
> Source: [.NET Microservices](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/)

A trimmed down version of the above diagram is shown below:

![pubsub-basic](./publish-subscribe-basics.png)
> Source: [.NET Microservices](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/)

## Installation

```bash
go get github.com/10hourlabs/whisper
```

## Usage

### Google Pub/Sub

```go
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
```
