package mqi

import (
	"log"
	"testing"

	"github.com/streadway/amqp"
)

func mockFunc(msg amqp.Delivery) error {
	log.Println("Inside of mockFunc | Test consumer function")
	log.Printf("%s\n", msg.Body)
	return nil
}

func getMockChannel() {
	GetChannel().UpdateChan() <- NewChannel(istore).
		WithBroadcast(NewBroadcast("test.broadcast").
			AddQueue(NewBroadcastQueue().
				AddTopic(NewTopic("test.bctp")).
				AddConsumer(NewConsumer("test.bccsm").WithFunc(mockFunc)))).
		WithExchange(NewExchange("test.exchange").
			AddQueue(NewQueue("test.q1").
				AddTopic(NewTopic("test.topic1")).
				AddConsumer(NewConsumer("test.consumer1").WithFunc(mockFunc))).
			AddQueue(NewQueue("test.q2").
				AddTopic(NewTopic("test.topic2")).
				AddConsumer(NewConsumer("test.consumer2").WithFunc(mockFunc))))
}

func getMockChannelWithOptionalQueue() {
	GetChannel().UpdateChan() <- NewChannel(istore).WithOptionalQueue(true).
		WithBroadcast(NewBroadcast("test.broadcast").
			AddQueue(NewBroadcastQueue().
				AddTopic(NewTopic("test.bctp")).
				AddConsumer(NewConsumer("test.bccsm").WithFunc(mockFunc)))).
		WithExchange(NewExchange("test.exchange.api"))
}

// Integration test
func TestSetup(t *testing.T) {
	getMockChannel()
	establish()
	setup()

	ch := GetChannel()
	if ch.Exchange().QueueAt(0).Ref() == nil {
		t.Error("Failed to setup")
	}
}

// Integration test
func TestAddTempQueue(t *testing.T) {
	getMockChannelWithOptionalQueue()
	establish()
	setup()
	q := NewTempQueue("test.q1.temp").
		AddTopic(NewTopic("test.topic1.temp")).
		AddConsumer(NewTempConsumer("test.consumer1.temp").WithFunc(mockFunc))
	q = AddTempQueue(q)

	if q.Ref() == nil {
		t.Error("Failed to setup")
	}
}
