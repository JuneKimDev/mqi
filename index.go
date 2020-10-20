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
	istore = InitStore()
}

// Get returns Store
func Get() Store { return istore }

// Start connects to RabbitMQ
func (ch channel) Start() {
	// Check if channel is ready
	if err := ch.isReady(); err != nil {
		log.Println("Please set all parts of Channel appropriately as a topic broker")
		return
	}

	log.Println("Connecting RabbitMQ sequence initialized")

	go ch.connect()
	ch.ErrChan() <- amqp.ErrClosed // Starts the sequence by passing an error

	// Wait until getting connected
	for !Get().GetChannel().IsStarted() {
	}
}

// Publish publishes a message with a topic to Exchange
func Publish(exchangeName string, topic string, msg amqp.Publishing) error {
	ch := Get().GetChannel()
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

// SendAck sends ACK to RabbitMQ
func SendAck(msg amqp.Delivery) {
	err := msg.Ack(false)
	if err != nil {
		log.Fatalf("Failed to send ACK to %d:%v", msg.DeliveryTag, err)
	}
}

// Close closes connection and channels
func Close() {
	ch := Get().GetChannel()
	ch.Conn().Close()
	ch.Sub().Close()
	ch.Pub().Close()
}
