# go-rmq

A simple and easy-to-use Go RabbitMQ client library built on top of the official AMQP library.

## Features

- Simple connection management
- Easy message publishing
- Message consuming with custom handlers
- Automatic queue declaration
- Persistent message delivery
- Proper error handling and logging

## Installation

```bash
go get github.com/idprm/go-rmq
```

## Quick Start

### Publishing Messages

```go
package main

import (
    "context"
    "log"

    "github.com/idprm/go-rmq"
)

func main() {
    // Create client
    client, err := rmq.NewClient(rmq.Config{
        URL: "amqp://guest:guest@localhost:5672/",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Publish message
    ctx := context.Background()
    err = client.PublishMessage(ctx, "my_queue", []byte("Hello, RabbitMQ!"))
    if err != nil {
        log.Fatal(err)
    }
}
```

### Consuming Messages

```go
package main

import (
    "fmt"
    "log"

    "github.com/idprm/go-rmq"
)

func main() {
    // Create client
    client, err := rmq.NewClient(rmq.Config{
        URL: "amqp://guest:guest@localhost:5672/",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Define message handler
    handler := func(message []byte) error {
        fmt.Printf("Received: %s\n", string(message))
        return nil
    }

    // Start consuming (blocking)
    err = client.ConsumeMessages("my_queue", handler)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

The `Config` struct supports the following options:

- `URL`: RabbitMQ connection URL (default: "amqp://guest:guest@localhost:5672/")

## API Reference

### Client Methods

- `NewClient(config Config) (*Client, error)` - Create a new RabbitMQ client
- `Close() error` - Close the connection
- `DeclareQueue(name string) error` - Declare a queue
- `PublishMessage(ctx context.Context, queueName string, message []byte) error` - Publish a message
- `ConsumeMessages(queueName string, handler func([]byte) error) error` - Consume messages
- `IsConnected() bool` - Check connection status

## Examples

See the `example/` directory for complete working examples:
- `example/publisher/` - Message publisher example
- `example/consumer/` - Message consumer example

## Requirements

- Go 1.18+
- RabbitMQ server

## License

MIT License
