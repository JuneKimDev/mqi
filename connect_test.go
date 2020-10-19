package mqi

import (
	"log"
	"testing"
)

func TestIsReadyFalse(t *testing.T) {
	falseInputs := make([]Channel, 8)
	// No Exchange
	falseInputs[0] = Get().GetChannel().WithURI("")
	// No Queue
	falseInputs[1] = Get().GetChannel().WithURI("").WithExchange(NewExchange(""))
	// No Topic
	falseInputs[2] = Get().GetChannel().WithURI("").WithExchange(NewExchange("").AddQueue(NewQueue("")))
	// No Consumer
	falseInputs[3] = Get().GetChannel().WithURI("").WithExchange(NewExchange("").AddQueue(NewQueue("").AddTopic(NewTopic(""))))
	// No ConsumerFunc
	falseInputs[4] = Get().GetChannel().WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer(""))))
	// No Topic in 2nd Queue
	falseInputs[5] = Get().GetChannel().WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer("").WithFunc(mockFunc))).
		AddQueue(NewQueue("")))
	// No Consumer in 2nd Queue
	falseInputs[6] = Get().GetChannel().WithURI("").WithExchange(NewExchange("").
		AddQueue(NewQueue("").AddTopic(NewTopic("")).AddConsumer(NewConsumer("").WithFunc(mockFunc))).
		AddQueue(NewQueue("").AddTopic(NewTopic(""))))
	// No ConsumerFunc in 2nd Queue
	falseInputs[7] = Get().GetChannel().WithURI("").WithExchange(NewExchange("").
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
	ch := getMockChannel()

	result := ch.isReady()
	if result != nil {
		t.Error("Supposed to pass but failed at channel readiness check")
	}
}

// Integration test
// func TestKeepTryingToConnectTrue(t *testing.T) {
// 	ch := Get().GetChannel()
// 	conn := keepTryingToConnect(ch.URI())
// 	log.Printf("%#v", conn)
// }

// Integration test
// This never stops
// func TestKeepTryingToConnectFalse(t *testing.T) {
// 	ch := Get().GetChannel().WithURI("amqp://wrong:addr@whatever.com:1234/")
// 	conn := keepTryingToConnect(ch.URI())
// 	log.Printf("%#v", conn)
// }

// Integration test
// func TestEstablish(t *testing.T) {
// 	ch := Get().GetChannel()
// 	result := ch.establish().(Channel)
// 	if result.Conn() == nil {
// 		t.Error("Failed to establish connection")
// 	}
// 	if result.Sub() == nil {
// 		t.Error("Failed to establish Sub channel")
// 	}
// 	if result.Pub() == nil {
// 		t.Error("Failed to establish Pub channel")
// 	}
// }

// Integration test
// func TestConnect(t *testing.T) {
// 	ch := getMockChannel()
//
// 	go ch.connect()
// 	ch.ErrChan() <- amqp.ErrClosed // Starts the sequence by passing an error
// 	ch = <-ch.ConnChan()
// 	log.Printf("%#v", ch)
// }
