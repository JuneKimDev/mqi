// mqi is RabbitMQ interface app for perme micro-service architecture
package mqi

import (
	"log"

	"github.com/streadway/amqp"
)

func mockFunc(ch Channel, msgs <-chan amqp.Delivery, kill <-chan bool) {
	log.Println("Inside of mockFunc | Test consumer function")
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

func mockFuncForever(ch Channel, msgs <-chan amqp.Delivery, kill <-chan bool) {
	log.Println("Inside of mockFuncForever | Test consumer function")
	for {
		select {
		case <-kill:
			log.Println("Shutting down a consumer")
			return
		case msg := <-msgs:

			log.Printf("MSG[%s]: %s\n", msg.RoutingKey, msg.Body)
			SendAck(msg)
		}
	}
}

func getMockChannelWithForverConsumer() Channel {
	return Get().GetChannel().
		WithExchange(NewExchange("test.exchange").
			AddQueue(NewQueue("test.q1").
				AddTopic(NewTopic("test.topic1")).
				AddConsumer(NewConsumer("test.consumer1").WithFunc(mockFuncForever))).
			AddQueue(NewQueue("test.q2").
				AddTopic(NewTopic("test.topic2")).
				AddConsumer(NewConsumer("test.consumer2").WithFunc(mockFuncForever))))
}

// Integration test
// func TestDisclosed(t *testing.T) {
// 	getMockChannelWithForverConsumer().Start()

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
// 	getMockChannelWithForverConsumer().Start()
// 	log.Println("The app is running in forever loop")
// 	forever.Wait() // Prevents service from exiting
// }
