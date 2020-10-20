package mqi

import "github.com/streadway/amqp"

// Consumer struct
type consumer struct {
	name string
	fn   ConsumerFunc
}

// ConsumerFunc Consumer worker function
type ConsumerFunc func(msgs <-chan amqp.Delivery, kill <-chan bool)

// Consumer interface
type Consumer interface {
	Name() string
	Func() ConsumerFunc
	WithName(name string) Consumer
	WithFunc(fn ConsumerFunc) Consumer
}

// NewConsumer constructor
func NewConsumer(name string) Consumer { return consumer{}.WithName(name) }

// Getters and withers
func (csm consumer) Name() string       { return csm.name }
func (csm consumer) Func() ConsumerFunc { return csm.fn }
func (csm consumer) WithName(name string) Consumer {
	csm.name = name
	return csm
}
func (csm consumer) WithFunc(fn ConsumerFunc) Consumer {
	csm.fn = fn
	return csm
}
