package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// Publish messages
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Hello RabbitMQ! Message #%d", i+1)

		err := client.PublishMessage(ctx, queueName, []byte(message))
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
			continue
		}

		fmt.Printf("Published: %s\n", message)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All messages published successfully!")
}
