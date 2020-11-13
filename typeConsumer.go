package mqi

import (
	"log"

	"github.com/streadway/amqp"
)

// Consumer struct
type consumer struct {
	name        string                    // Name of Consumer
	typeString  string                    // type
	isAutoAck   bool                      // Auto-ack
	isExclusive bool                      // Exclusive
	isNoLocal   bool                      // No-local
	isNoWait    bool                      // No-wait
	fn          func(amqp.Delivery) error // consumer function
}

// Consumer interface
type Consumer interface {
	Name() string
	TypeString() string
	IsAutoAck() bool
	IsExclusive() bool
	IsNoLocal() bool
	IsNoWait() bool
	Func() func(msg amqp.Delivery) error
	WithName(name string) Consumer
	WithTypeString(str string) Consumer
	WithAutoAck(b bool) Consumer
	WithExclusive(b bool) Consumer
	WithNoLocal(b bool) Consumer
	WithNoWait(b bool) Consumer
	WithFunc(fn func(msg amqp.Delivery) error) Consumer
	bindWith(q Queue)
}

// NewConsumer constructs a consumer
func NewConsumer(name string) Consumer { return consumer{}.WithName(name) }

// NewTempConsumer constructs a temporary consumer
func NewTempConsumer(name string) Consumer { return consumer{}.WithName(name).WithTypeString("temp") }

// Getters and withers
func (csm consumer) Name() string                        { return csm.name }
func (csm consumer) TypeString() string                  { return csm.typeString }
func (csm consumer) IsAutoAck() bool                     { return csm.isAutoAck }
func (csm consumer) IsExclusive() bool                   { return csm.isExclusive }
func (csm consumer) IsNoLocal() bool                     { return csm.isNoLocal }
func (csm consumer) IsNoWait() bool                      { return csm.isNoWait }
func (csm consumer) Func() func(msg amqp.Delivery) error { return csm.fn }
func (csm consumer) WithName(name string) Consumer {
	csm.name = name
	return csm
}
func (csm consumer) WithTypeString(str string) Consumer {
	csm.typeString = str
	return csm
}
func (csm consumer) WithAutoAck(b bool) Consumer {
	csm.isAutoAck = b
	return csm
}
func (csm consumer) WithExclusive(b bool) Consumer {
	csm.isExclusive = b
	return csm
}
func (csm consumer) WithNoLocal(b bool) Consumer {
	csm.isNoLocal = b
	return csm
}
func (csm consumer) WithNoWait(b bool) Consumer {
	csm.isNoWait = b
	return csm
}
func (csm consumer) WithFunc(fn func(msg amqp.Delivery) error) Consumer {
	csm.fn = fn
	return csm
}

// runs in goroutine
func (csm consumer) bindWith(q Queue) {
	ch := GetChannel()
	qName := q.Name()

	msgs, err := ch.Sub().Consume(
		qName,           // queue name
		csm.name,        // consumer name
		csm.isAutoAck,   // auto-ack
		csm.isExclusive, // exclusive
		csm.isNoLocal,   // no-local
		csm.isNoWait,    // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatalf("Failed to consume from [%s]: %v", qName, err)
		// Exit b/c failed to setup
	}
	log.Printf("Consuming a queue [%s]...\n", qName)

	if csm.TypeString() == "temp" {
		// Temporary consumer doesn't loop forever
		// Temporary consumer will be used only once and will delete the queue
		go func() {
			msg := <-msgs
			if err := csm.Func()(msg); err != nil {
				log.Printf("Error in consumer function: %v\n", err)
				sendNack(msg)
			} else {
				sendAck(msg)
			}

			// Delete the Queue
			GetChannel().Sub().QueueDelete(
				qName, // Queue name
				false, // ifUnused
				false, // ifEmpty
				false, // noWait
			)
			log.Printf("[%s] Queue deleted\n", qName)
		}()
	} else {
		// Normal consumer loops forever until kill signal comes in
		go func() {
			for {
				select {
				case <-ch.KillChan():
					log.Println("Shutting down a consumer")
					return
				case msg := <-msgs:
					if err := csm.Func()(msg); err != nil {
						log.Printf("Error in consumer function: %v\n", err)
						sendNack(msg)
					} else {
						sendAck(msg)
					}
				}
			}
		}()
	}
}

// sendAck sends ACK to RabbitMQ
func sendAck(msg amqp.Delivery) {
	err := msg.Ack(false)
	if err != nil {
		log.Printf("Failed to send ACK to %d:%v", msg.DeliveryTag, err)
	}
}

// sendNack sends NACK to RabbitMQ
func sendNack(msg amqp.Delivery) {
	err := msg.Nack(false, true)
	if err != nil {
		log.Printf("Failed to send NACK to %d:%v", msg.DeliveryTag, err)
	}
}
