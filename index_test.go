package mqi

import (
	"log"

	"github.com/streadway/amqp"
)

func mockFunc(msg amqp.Delivery) error {
	log.Println("Inside of mockFunc | Test consumer function")
	return nil
}

func getMockChannel() Channel {
	return Get().GetChannel().
		WithExchange(NewExchange("test.exchange").
			AddQueue(NewQueue("test.q1").
				AddTopic(NewTopic("test.topic1")).
				AddConsumer(NewConsumer("test.consumer1").WithFunc(mockFunc))).
			AddQueue(NewQueue("test.q2").
				AddTopic(NewTopic("test.topic2")).
				AddConsumer(NewConsumer("test.consumer2").WithFunc(mockFunc))))
}

// Integration test
// func TestDisclosed(t *testing.T) {
// 	getMockChannel().Start()

// 	// Send a publish
// 	msg := amqp.Publishing{ContentType: "text/plain", Body: []byte("Hello From Test Publisher")}
// 	err := Publish("test.exchange", "test.topic1", msg)
// 	if err != nil {
// 		t.Errorf("Failed to publish: %v", err)
// 	}

// 	// Send a close event to watch reconnection in action
// 	Get().GetChannel().ErrChan() <- amqp.ErrClosed
// 	Publish("test.exchange", "test.topic1", msg)
// 	if err != nil {
// 		t.Errorf("Failed to publish: %v", err)
// 	}
// 	time.Sleep(time.Second * 5)
// }

// Integration test
// func TestMain(t *testing.T) {
// 	var forever sync.WaitGroup
// 	forever.Add(1)
// 	getMockChannel().Start()
// 	log.Println("The app is running in forever loop")
// 	forever.Wait() // Prevents service from exiting
// }
