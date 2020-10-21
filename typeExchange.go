package mqi

// Exchange struct
type exchange struct {
	name string  // Name of Exchange
	qs   []Queue // Declared queues
}

// Exchange interface
type Exchange interface {
	Name() string
	CountQueues() int
	QueueAt(i int) Queue
	CountAllConsumers() int
	WithName(name string) Exchange
	AddQueue(q Queue) Exchange
	UpdateQueue(q Queue) Exchange
}

// NewExchange Exchange Constructor
func NewExchange(name string) Exchange { return exchange{}.WithName(name) }

// Getters and withers
func (ex exchange) Name() string        { return ex.name }
func (ex exchange) CountQueues() int    { return len(ex.qs) }
func (ex exchange) QueueAt(i int) Queue { return ex.qs[i] }
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
		}
	}
	return ex
}
