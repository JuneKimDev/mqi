// Package mqi is Message Queue Interface using AMQP
//
// Requires active running instance of RabbitMQ server
package mqi

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
)

var istore Store

func init() {
	istore = initStore()
}

// GetChannel returns current Channel
func GetChannel() Channel {
	resChan := make(chan Channel)
	istore.ReqChan() <- resChan
	return <-resChan
}

// Start connects to RabbitMQ
func (ch channel) Start() {
	// Check if channel is ready
	if err := ch.isReady(); err != nil {
		log.Println("Please set all parts of Channel appropriately as a topic broker")
		return
	}

	log.Println("Connecting RabbitMQ sequence initialized")

	ch.UpdateChan() <- ch // Updates store
	go connect()
	ch.ErrChan() <- amqp.ErrClosed // Starts the sequence by passing an error

	// Wait until getting connected
	for !GetChannel().IsStarted() {
	}
}

// AddTempQueue adds a queue which gets auto-deleted after execution of consumer function
func AddTempQueue(q Queue) Queue {
	log.Println("Adding a queue...")
	q = declareTempQueue(q)

	for j := 0; j < q.CountTopics(); j++ {
		tp := q.TopicAt(j)
		bindQueueWith(q, tp)
	}

	for j := 0; j < q.CountConsumers(); j++ {
		csm := q.ConsumerAt(j)
		bindTempConsumerWith(q, csm)
	}
	return q
}

// Publish publishes a message with a topic to Exchange
func Publish(exchangeName string, topic string, msg amqp.Publishing) error {
	ch := GetChannel()
	if !ch.IsStarted() {
		return errors.New("RabbitMQ is not connected currently")
	}
	err := ch.Pub().Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		msg,          // message
	)
	if err != nil {
		log.Fatalf("Failed to Publish [%s]: %v", topic, err)
	}

	log.Printf("Sent a msg (%s) successfully to [%s]\n", topic, exchangeName)
	return nil
}

// Close closes connection and channels
func Close() {
	ch := GetChannel()
	ch.Conn().Close()
	ch.Sub().Close()
	ch.Pub().Close()
}
