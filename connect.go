package mqi

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func (ch channel) isReady() error {
	if ch.URI() == "amqp://guest:guest@localhost:5672/" {
		log.Println("Warning! You are using the default connection string for RabbitMQ.")
	}
	if ch.Exchange() == nil {
		return errors.New("No Exchange is set in Channel")
	}

	// API gateway doesn't require queue declaration at the time of setup
	// But other services need full queue/topic/consumer declaration at the time of setup
	if !ch.IsOptionalQueue() && ch.Exchange().CountQueues() == 0 {
		return fmt.Errorf("No Queues are set in Exchange %s", ch.Exchange().Name())
	}
	for i := 0; i < ch.Exchange().CountQueues(); i++ {
		if ch.Exchange().QueueAt(i).CountTopics() == 0 {
			return fmt.Errorf("No Topics are set in Queue %s", ch.Exchange().QueueAt(i).Name())
		}
		if ch.Exchange().QueueAt(i).CountConsumers() == 0 {
			return fmt.Errorf("No Consumers are set in Queue %s", ch.Exchange().QueueAt(i).Name())
		}
		for j := 0; j < ch.Exchange().QueueAt(i).CountConsumers(); j++ {
			if ch.Exchange().QueueAt(i).ConsumerAt(j).Func() == nil {
				return fmt.Errorf("No ConsumerFunc is set in Consumer %s of Queue %s",
					ch.Exchange().QueueAt(i).ConsumerAt(j).Name(), ch.Exchange().QueueAt(i).Name())
			}
		}
	}
	return nil
}

// connect (runs forever in goroutine)
func connect() {

	// Receive trigger
	ch := GetChannel()
	for trigger := range ch.ErrChan() {
		if trigger == amqp.ErrClosed { // Connection has been closed
			log.Println("Disconnected from RabbitMQ")

			// IsStarted has to be check with current channel
			ch = GetChannel()
			if ch.Exchange().CountQueues() != 0 && ch.IsStarted() {
				// Terminate all consumer goroutines
				log.Println("Sending kill-signals to Consumer goroutines")
				for i := 0; i < ch.Exchange().CountAllConsumers(); i++ {
					ch.KillChan() <- true
				}
			}

			// Set isStarted false
			updateChannelStarted(false)

			// Connect to RabbitMQ
			log.Println("Connecting to RabbitMQ...")
			establish()
			setup()
			updateChannelStarted(true)
			log.Println("Successfully connected to RabbitMQ")
		} else {
			log.Fatalf("Unexpected error occurred: %v", trigger) // Exit b/c failed to setup
		}
	}
}

func establish() {
	// Establish connection
	ch := GetChannel()
	conn := keepTryingToConnect(ch.URI())
	conn.NotifyClose(ch.ErrChan()) // Emits close event to channel

	// Establish subscription channel
	subCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a sub channel: %v", err) // Exit b/c failed to setup
	}
	subCh.Qos(
		ch.Prefetch(), // prefetch count
		0,             // prefetch size
		false,         // global
	)
	log.Println("Sub channel established successfully")

	// Establish publish channel
	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a pub channel: %v", err) // Exit b/c failed to setup
	}
	pubCh.Qos(
		ch.Prefetch(), // prefetch count
		0,             // prefetch size
		false,         // global
	)

	log.Println("Pub channel established successfully")

	ch.UpdateChan() <- ch.WithConn(conn).WithSub(subCh).WithPub(pubCh)
}

func keepTryingToConnect(uri string) *amqp.Connection {
	var serverAt string
	uriParse := strings.Split(uri, "@")

	if len(uriParse) != 2 {
		log.Fatalf("Connection string is malformatted: want 2, got %d", len(uriParse)) // Exit b/c failed to setup
	} else {
		serverAt = uriParse[1]
	}

	var delay int
	var delayCount int
	// Loop until got connection
	for {
		conn, err := amqp.Dial(uri)
		if err == nil {
			log.Printf("Connected to RabbitMQ server at %s successfully\n", serverAt)
			delay, delayCount = 0, 0
			return conn
		}

		log.Println(err)

		// Delay timer in stepping increment mode
		// Every 3 retry, delay time increases 500ms
		time.Sleep(time.Duration(delay) * time.Millisecond)
		delayCount++
		if delayCount == 3 {
			delayCount = 0
			delay += 500
			// Breaker for unit test
			if uri == "amqp://wrong:addr@test.debug:1234/" {
				break
			}
		}
		log.Printf("Trying to reconnect to RabbitMQ at %s\n", serverAt)
	}
	return nil
}

func updateChannelStarted(b bool) {
	ch := GetChannel()
	ch.UpdateChan() <- ch.WithStarted(b)
}
