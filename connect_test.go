package mqi

import (
	"log"
	"testing"

	"github.com/streadway/amqp"
)

func TestIsReadyFalse(t *testing.T) {
	falseInputs := make([]Channel, 8)
	// No Exchange
	falseInputs[0] = NewChannel(istore).WithURI("")
	// No Queue
	falseInputs[1] = NewChannel(istore).WithURI("").WithExchange(NewExchange(""))
	// No Topic
	falseInputs[2] = NewChannel(istore).WithURI("").WithExchange(NewExchange("").AddQueue(NewQueue("")))
	// No Consumer
	falseInputs[3] = NewChannel(istore).WithURI("").WithExchange(NewExchange("").AddQueue(NewQueue("").AddTopic(NewTopic(""))))
	// No ConsumerFunc
	falseInputs[4] = NewChannel(istore).WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer(""))))
	// No Topic in 2nd Queue
	falseInputs[5] = NewChannel(istore).WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer("").WithFunc(mockFunc))).
		AddQueue(NewQueue("")))
	// No Consumer in 2nd Queue
	falseInputs[6] = NewChannel(istore).WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer("").WithFunc(mockFunc))).
		AddQueue(NewQueue("").AddTopic(NewTopic(""))))
	// No ConsumerFunc in 2nd Queue
	falseInputs[7] = NewChannel(istore).WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer("").WithFunc(mockFunc))).
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer(""))))

	for i, in := range falseInputs {
		result := in.isReady()
		if result == nil {
			t.Errorf("Readiness check at (%d) outputs false positive", i)
		} else {
			log.Println(result)
		}
	}
}

func TestIsReadyTrue(t *testing.T) {
	getMockChannel()

	result := GetChannel().isReady()
	if result != nil {
		t.Error("Supposed to pass but failed at channel readiness check")
	}
}

func TestIsReadyWithOptionalQueue(t *testing.T) {
	getMockChannelWithOptionalQueue()

	result := GetChannel().isReady()
	if result != nil {
		t.Error("Supposed to pass but failed at channel readiness check")
	}
}

// Integration test
func TestKeepTryingToConnectTrue(t *testing.T) {
	getMockChannel()
	ch := GetChannel()
	conn := keepTryingToConnect(ch.URI())
	log.Printf("%#v", conn)
}

// Integration test (This will cause test to never stop)
func TestKeepTryingToConnectFalse(t *testing.T) {
	getMockChannel()
	ch := GetChannel().WithURI("amqp://wrong:addr@test.debug:1234/")
	conn := keepTryingToConnect(ch.URI())
	if conn != nil {
		t.Errorf("Supposed to fail but got %v", conn)
	}
}

// Integration test
func TestEstablish(t *testing.T) {
	getMockChannel()
	establish()
	ch := GetChannel()
	if ch.Conn() == nil {
		t.Error("Failed to establish connection")
	}
	if ch.Sub() == nil {
		t.Error("Failed to establish Sub channel")
	}
	if ch.Pub() == nil {
		t.Error("Failed to establish Pub channel")
	}
}

// Integration test
func TestConnect(t *testing.T) {
	getMockChannel()
	go connect()
	GetChannel().ErrChan() <- amqp.ErrClosed // Starts the sequence by passing an error
	for !GetChannel().IsStarted() {
	}
	log.Printf("%#v", GetChannel())
}
