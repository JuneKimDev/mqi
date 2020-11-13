package mqi

import (
	"log"

	"github.com/streadway/amqp"
)

// Queue struc
type queue struct {
	name                string      // Name of Queue
	typeString          string      // type
	isDurable           bool        // Duraility
	isAutoDeleteEnabled bool        // Auto-delete enabled
	isExclusive         bool        // Exclusivity
	isNoWait            bool        // No-wait
	args                amqp.Table  // arguments
	ref                 *amqp.Queue // Reference to the attached AMQP Queue
	tps                 []Topic     // Binded topics
	csms                []Consumer  // Binded consumers
}

//Queue interface
type Queue interface {
	Name() string
	TypeString() string
	IsDurable() bool
	IsAutoDeleteEnabled() bool
	IsExclusive() bool
	IsNoWait() bool
	Args() amqp.Table
	Ref() *amqp.Queue
	CountTopics() int
	TopicAt(i int) Topic
	CountConsumers() int
	ConsumerAt(i int) Consumer
	WithName(name string) Queue
	withTypeString(str string) Queue
	WithDurable(b bool) Queue
	WithAutoDeleteEnabled(b bool) Queue
	WithExclusive(b bool) Queue
	WithNoWait(b bool) Queue
	WithArgs(t amqp.Table) Queue
	WithRef(ref *amqp.Queue) Queue
	AddTopic(tp Topic) Queue
	AddConsumer(cs Consumer) Queue
	declare() Queue
}

// NewQueue constructs a queue
func NewQueue(name string) Queue { return queue{}.WithName(name).WithDurable(true) }

// NewTempQueue constructs a temporary queue
func NewTempQueue(name string) Queue {
	return queue{}.WithName(name).withTypeString("temp").WithDurable(true).WithAutoDeleteEnabled(true)
}

// NewBroadcastQueue constructs a broadcast queue
func NewBroadcastQueue() Queue {
	return queue{}.withTypeString("cast").WithAutoDeleteEnabled(true).WithExclusive(true)
}

// Getters and withers
func (q queue) Name() string              { return q.name }
func (q queue) TypeString() string        { return q.typeString }
func (q queue) IsDurable() bool           { return q.isDurable }
func (q queue) IsAutoDeleteEnabled() bool { return q.isAutoDeleteEnabled }
func (q queue) IsExclusive() bool         { return q.isExclusive }
func (q queue) IsNoWait() bool            { return q.isNoWait }
func (q queue) Args() amqp.Table          { return q.args }
func (q queue) Ref() *amqp.Queue          { return q.ref }
func (q queue) CountTopics() int          { return len(q.tps) }
func (q queue) TopicAt(i int) Topic       { return q.tps[i] }
func (q queue) CountConsumers() int       { return len(q.csms) }
func (q queue) ConsumerAt(i int) Consumer { return q.csms[i] }
func (q queue) WithName(name string) Queue {
	q.name = name
	return q
}
func (q queue) withTypeString(str string) Queue {
	q.typeString = str
	return q
}
func (q queue) WithDurable(b bool) Queue {
	q.isDurable = b
	return q
}
func (q queue) WithAutoDeleteEnabled(b bool) Queue {
	q.isAutoDeleteEnabled = b
	return q
}
func (q queue) WithExclusive(b bool) Queue {
	q.isExclusive = b
	return q
}
func (q queue) WithNoWait(b bool) Queue {
	q.isNoWait = b
	return q
}
func (q queue) WithArgs(t amqp.Table) Queue {
	q.args = t
	return q
}
func (q queue) WithRef(ref *amqp.Queue) Queue {
	q.ref = ref
	return q
}
func (q queue) AddTopic(tp Topic) Queue {
	q.tps = append(q.tps, tp)
	return q
}
func (q queue) AddConsumer(csm Consumer) Queue {
	q.csms = append(q.csms, csm)
	return q
}

func (q queue) declare() Queue {
	ch := GetChannel()
	qName := q.name
	ref, err := ch.Sub().QueueDeclare(
		qName,                 // queue name
		q.isDurable,           // durable; it must follow durability of exchange with which it is binded
		q.isAutoDeleteEnabled, // auto-delete when there are no remaining consumers
		q.isExclusive,         // exclusive
		q.isNoWait,            // no-wait
		q.args,                // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue [%s]: %v", qName, err)
		// Exit b/c failed to setup
	}
	dq := q.WithRef(&ref)
	if dq.TypeString() == "cast" {
		ch.UpdateChan() <- ch.WithBroadcast(ch.Broadcast().UpdateQueue(dq))
	} else if dq.TypeString() != "temp" {
		ch.UpdateChan() <- ch.WithExchange(ch.Exchange().UpdateQueue(dq))
	}
	log.Printf("Declared a queue [%s] successfully\n", qName)
	return dq
}
