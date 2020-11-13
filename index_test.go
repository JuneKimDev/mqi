package mqi

import (
	"testing"

	"github.com/streadway/amqp"
)

// Integration test
func TestMQIFlow(t *testing.T) {
	getMockChannel()
	GetChannel().Start()

	// Send a publish
	msg := amqp.Publishing{ContentType: "text/plain", Body: []byte("Hello From Test Publisher")}
	err := Publish("test.exchange", "test.topic1", msg)
	if err != nil {
		t.Errorf("Failed to publish: %v", err)
	}
	err = Publish("test.exchange", "test.topic1", msg)
	if err != nil {
		t.Errorf("Failed to publish: %v", err)
	}
	// Close()
}
