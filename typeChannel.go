package mqi

import "github.com/streadway/amqp"

// Channel struct
type channel struct {
	uri             string           // RabbitMQ connection string
	prefetch        int              // Number of prefetch
	ex              Exchange         // Declared exchange
	cast            Exchange         // Broadcasting exchange
	conn            *amqp.Connection // RabbitMQ connection
	sub             *amqp.Channel    // RabbitMQ channel for subscription
	pub             *amqp.Channel    // RabbitMQ channel for publishing
	updateChan      chan<- Channel   // Gochannel for updating the state
	killChan        chan bool        // Gochannel to send kill signal to consumers with forever loop
	errChan         chan *amqp.Error // Gochannel to report disconnection to RabbitMQ server
	isStarted       bool             // Start indicator
	isOptionalQueue bool             // Skips number of queue at the time of setup
}

// Channel interface
type Channel interface {
	URI() string
	Prefetch() int
	Exchange() Exchange
	Broadcast() Exchange
	Conn() *amqp.Connection
	Sub() *amqp.Channel
	Pub() *amqp.Channel
	UpdateChan() chan<- Channel
	KillChan() chan bool
	ErrChan() chan *amqp.Error
	IsStarted() bool
	IsOptionalQueue() bool
	WithURI(uri string) Channel
	WithPrefetch(prefetch int) Channel
	WithExchange(ex Exchange) Channel
	WithBroadcast(ex Exchange) Channel
	WithConn(conn *amqp.Connection) Channel
	WithSub(sub *amqp.Channel) Channel
	WithPub(pub *amqp.Channel) Channel
	WithUpdateChan(c chan<- Channel) Channel
	WithKillChan(c chan bool) Channel
	WithErrChan(c chan *amqp.Error) Channel
	WithStarted(b bool) Channel
	WithOptionalQueue(b bool) Channel
	Start()
	isReady() error
}

// NewChannel constructor
func NewChannel(st Store) Channel {
	return channel{}.
		WithURI("amqp://guest:guest@localhost:5672/").
		WithPrefetch(100).
		WithUpdateChan(st.UpdateChan()).
		WithKillChan(make(chan bool)).
		WithErrChan(make(chan *amqp.Error))
}

// Getters and withers
func (ch channel) URI() string                { return ch.uri }
func (ch channel) Prefetch() int              { return ch.prefetch }
func (ch channel) Exchange() Exchange         { return ch.ex }
func (ch channel) Broadcast() Exchange        { return ch.cast }
func (ch channel) Conn() *amqp.Connection     { return ch.conn }
func (ch channel) Sub() *amqp.Channel         { return ch.sub }
func (ch channel) Pub() *amqp.Channel         { return ch.pub }
func (ch channel) UpdateChan() chan<- Channel { return ch.updateChan }
func (ch channel) KillChan() chan bool        { return ch.killChan }
func (ch channel) ErrChan() chan *amqp.Error  { return ch.errChan }
func (ch channel) IsStarted() bool            { return ch.isStarted }
func (ch channel) IsOptionalQueue() bool      { return ch.isOptionalQueue }
func (ch channel) WithURI(uri string) Channel {
	ch.uri = uri
	return ch
}
func (ch channel) WithPrefetch(prefetch int) Channel {
	ch.prefetch = prefetch
	return ch
}
func (ch channel) WithExchange(ex Exchange) Channel {
	ch.ex = ex
	return ch
}
func (ch channel) WithBroadcast(cast Exchange) Channel {
	ch.cast = cast
	return ch
}
func (ch channel) WithConn(conn *amqp.Connection) Channel {
	ch.conn = conn
	return ch
}
func (ch channel) WithSub(sub *amqp.Channel) Channel {
	ch.sub = sub
	return ch
}
func (ch channel) WithPub(pub *amqp.Channel) Channel {
	ch.pub = pub
	return ch
}
func (ch channel) WithUpdateChan(c chan<- Channel) Channel {
	ch.updateChan = c
	return ch
}
func (ch channel) WithKillChan(c chan bool) Channel {
	ch.killChan = c
	return ch
}
func (ch channel) WithErrChan(c chan *amqp.Error) Channel {
	ch.errChan = c
	return ch
}
func (ch channel) WithStarted(b bool) Channel {
	ch.isStarted = b
	return ch
}
func (ch channel) WithOptionalQueue(b bool) Channel {
	ch.isOptionalQueue = b
	return ch
}
