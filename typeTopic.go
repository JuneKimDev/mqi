package mqi

import "log"

// Topic struct
type topic struct {
	name string
}

// Topic interface
type Topic interface {
	Name() string
	WithName(name string) Topic
	bindWith(q Queue)
}

// NewTopic constructor
func NewTopic(name string) Topic { return topic{}.WithName(name) }

// Name getter
func (tp topic) Name() string { return tp.name }

// WithName setter
func (tp topic) WithName(name string) Topic {
	tp.name = name
	return tp
}

func (tp topic) bindWith(q Queue) {
	ch := GetChannel()
	var exName string
	if q.TypeString() == "cast" {
		exName = ch.Broadcast().Name()
	} else {
		exName = ch.Exchange().Name()
	}
	qName := q.Name()
	tpName := tp.Name()
	err := ch.Sub().QueueBind(
		qName,  // queue name
		tpName, // routing key
		exName, // exchange name
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue [%s:%s] with [%s]: %v", qName, tpName, exName, err)
		// Exit b/c failed to setup
	}
	log.Printf("Binded a queue [%s:%s] with [%s] successfully\n", qName, tpName, exName)
}
