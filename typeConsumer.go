package mqi

import "github.com/streadway/amqp"

// Consumer struct
type consumer struct {
	name string
	fn   func(amqp.Delivery) error
}

// Consumer interface
type Consumer interface {
	Name() string
	Func() func(msg amqp.Delivery) error
	WithName(name string) Consumer
	WithFunc(fn func(msg amqp.Delivery) error) Consumer
}

// NewConsumer constructor
func NewConsumer(name string) Consumer { return consumer{}.WithName(name) }

// Getters and withers
func (csm consumer) Name() string                        { return csm.name }
func (csm consumer) Func() func(msg amqp.Delivery) error { return csm.fn }
func (csm consumer) WithName(name string) Consumer {
	csm.name = name
	return csm
}
func (csm consumer) WithFunc(fn func(msg amqp.Delivery) error) Consumer {
	csm.fn = fn
	return csm
}
