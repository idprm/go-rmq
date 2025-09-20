package main

import (
	"fmt"
	"log"

	"github.com/idprm/go-rmq"
)

func main() {
	// Create RabbitMQ client
	client, err := rmq.NewClient(rmq.Config{
		URL: "amqp://guest:guest@localhost:5672/",
	})
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	// Queue name
	queueName := "test_queue"

	// Message handler
	handler := func(message []byte) error {
		fmt.Printf("Received message: %s\n", string(message))
		return nil
	}

	// Start consuming messages
	fmt.Printf("Starting to consume messages from queue '%s'...\n", queueName)
	err = client.ConsumeMessages(queueName, handler)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}
