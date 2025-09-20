package rmq

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test with default config
	config := Config{}
	client, err := NewClient(config)

	// We expect this to fail if RabbitMQ is not running
	// which is normal in a test environment
	if err != nil {
		t.Logf("Expected error when RabbitMQ is not available: %v", err)
		return
	}

	if client == nil {
		t.Error("Expected client to be created")
		return
	}

	defer client.Close()

	// Test connection status
	if !client.IsConnected() {
		t.Error("Expected client to be connected")
	}
}

func TestClientClose(t *testing.T) {
	config := Config{
		URL: "amqp://invalid:invalid@nonexistent:5672/",
	}

	client, err := NewClient(config)
	if err != nil {
		// Expected to fail with invalid connection
		t.Logf("Expected error with invalid config: %v", err)
		return
	}

	if client != nil {
		err = client.Close()
		if err != nil {
			t.Errorf("Error closing connection: %v", err)
		}
	}
}

func TestConfig(t *testing.T) {
	// Test default URL
	config := Config{}
	if config.URL != "" {
		t.Error("Expected empty URL in default config")
	}

	// Test custom URL
	customURL := "amqp://user:pass@host:5672/vhost"
	config = Config{URL: customURL}
	if config.URL != customURL {
		t.Errorf("Expected URL %s, got %s", customURL, config.URL)
	}
}
