package rmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client represents a RabbitMQ client
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// Config holds configuration for RabbitMQ connection
type Config struct {
	URL string
}

// NewClient creates a new RabbitMQ client
func NewClient(config Config) (*Client, error) {
	if config.URL == "" {
		config.URL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: channel,
		url:     config.URL,
	}, nil
}

// Close closes the RabbitMQ connection
func (c *Client) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
			return err
		}
	}
	return nil
}

// DeclareQueue declares a queue
func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// PublishMessage publishes a message to a queue
func (c *Client) PublishMessage(ctx context.Context, queueName string, message []byte) error {
	// Ensure queue exists
	if err := c.DeclareQueue(queueName); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         message,
			DeliveryMode: amqp.Persistent, // make message persistent
			Timestamp:    time.Now(),
		},
	)
}

// ConsumeMessages consumes messages from a queue
func (c *Client) ConsumeMessages(queueName string, handler func([]byte) error) error {
	// Ensure queue exists
	if err := c.DeclareQueue(queueName); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				d.Nack(false, true) // negative acknowledge, requeue
			} else {
				d.Ack(false) // acknowledge message
			}
		}
	}()

	log.Printf("Waiting for messages from queue '%s'. To exit press CTRL+C", queueName)
	<-forever

	return nil
}

// IsConnected checks if the client is connected to RabbitMQ
func (c *Client) IsConnected() bool {
	return c.conn != nil && !c.conn.IsClosed()
}
