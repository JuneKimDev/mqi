package mqi

import (
	"log"

	"github.com/streadway/amqp"
)

// Exchange struct
type exchange struct {
	name                string     // Name of Exchange
	typeString          string     // Type of Exchange
	isDurable           bool       // Duraility
	isAutoDeleteEnabled bool       // Auto-delete enabled
	isExclusive         bool       // Exclusivity
	isInternal          bool       // Internal
	isNoWait            bool       // No-wait
	args                amqp.Table // arguments
	qs                  []Queue    // Declared queues
}

// Exchange interface
type Exchange interface {
	Name() string
	TypeString() string
	IsDurable() bool
	IsAutoDeleteEnabled() bool
	IsExclusive() bool
	IsInternal() bool
	IsNoWait() bool
	Args() amqp.Table
	CountQueues() int
	QueueAt(i int) Queue
	CountAllConsumers() int
	WithName(name string) Exchange
	WithTypeString(str string) Exchange
	WithDurable(b bool) Exchange
	WithAutoDeleteEnabled(b bool) Exchange
	WithExclusive(b bool) Exchange
	WithInternal(b bool) Exchange
	WithNoWait(b bool) Exchange
	WithArgs(t amqp.Table) Exchange
	AddQueue(q Queue) Exchange
	UpdateQueue(q Queue) Exchange
	declare()
}

// NewExchange constructs an exchange
func NewExchange(name string) Exchange {
	return exchange{}.
		WithName(name).
		WithTypeString("topic").
		WithDurable(true)
}

// NewBroadcast constructs an exchange for broadcast
func NewBroadcast(name string) Exchange {
	return exchange{}.
		WithName(name).
		WithTypeString("fanout")
}

// Getters and withers
func (ex exchange) Name() string              { return ex.name }
func (ex exchange) TypeString() string        { return ex.typeString }
func (ex exchange) IsDurable() bool           { return ex.isDurable }
func (ex exchange) IsAutoDeleteEnabled() bool { return ex.isAutoDeleteEnabled }
func (ex exchange) IsExclusive() bool         { return ex.isExclusive }
func (ex exchange) IsInternal() bool          { return ex.isInternal }
func (ex exchange) IsNoWait() bool            { return ex.isNoWait }
func (ex exchange) Args() amqp.Table          { return ex.args }
func (ex exchange) CountQueues() int          { return len(ex.qs) }
func (ex exchange) QueueAt(i int) Queue       { return ex.qs[i] }
func (ex exchange) CountAllConsumers() int {
	n := 0
	for i := 0; i < ex.CountQueues(); i++ {
		n += ex.QueueAt(i).CountConsumers()
	}
	return n
}
func (ex exchange) WithName(name string) Exchange {
	ex.name = name
	return ex
}
func (ex exchange) WithTypeString(str string) Exchange {
	ex.typeString = str
	return ex
}
func (ex exchange) WithDurable(b bool) Exchange {
	ex.isDurable = b
	return ex
}
func (ex exchange) WithAutoDeleteEnabled(b bool) Exchange {
	ex.isAutoDeleteEnabled = b
	return ex
}
func (ex exchange) WithExclusive(b bool) Exchange {
	ex.isExclusive = b
	return ex
}
func (ex exchange) WithInternal(b bool) Exchange {
	ex.isInternal = b
	return ex
}
func (ex exchange) WithNoWait(b bool) Exchange {
	ex.isNoWait = b
	return ex
}
func (ex exchange) WithArgs(t amqp.Table) Exchange {
	ex.args = t
	return ex
}

// AddQueue adds a Queue
func (ex exchange) AddQueue(q Queue) Exchange {
	ex.qs = append(ex.qs, q)
	return ex
}

// UpdateQueue updates a Queue
func (ex exchange) UpdateQueue(q Queue) Exchange {
	for i := 0; i < ex.CountQueues(); i++ {
		if ex.qs[i].Name() == q.Name() {
			ex.qs = append(append(ex.qs[:i], q), ex.qs[i+1:]...)
			break
		}
	}
	return ex
}

// Declare declares exchange in Sub channel
func (ex exchange) declare() {
	ch := GetChannel()
	err := ch.Sub().ExchangeDeclare(
		ex.name,                // name
		ex.typeString,          // type
		ex.isDurable,           // durable
		ex.isAutoDeleteEnabled, // auto-deleted
		ex.isInternal,          // internal
		ex.isNoWait,            // no-wait
		ex.args,                // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
		// Exit b/c failed to setup
	}
	log.Printf("Declared an exchange [%s] successfully\n", ex.name)
}
