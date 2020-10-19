package mqi

// Queue struc
type queue struct {
	name string     // Name of Queue
	tps  []Topic    // Binded topics
	csms []Consumer // Binded consumers
}

//Queue interface
type Queue interface {
	Name() string
	CountTopics() int
	TopicAt(i int) Topic
	CountConsumers() int
	ConsumerAt(i int) Consumer
	WithName(name string) Queue
	AddTopic(tp Topic) Queue
	AddConsumer(cs Consumer) Queue
}

// NewQueue constructor
func NewQueue(name string) Queue { return queue{}.WithName(name) }

// Name getter
func (q queue) Name() string { return q.name }

// CountTopics counts topics
func (q queue) CountTopics() int { return len(q.tps) }

// TopicAt gets at index
func (q queue) TopicAt(i int) Topic { return q.tps[i] }

// CountConsumers counts consumer
func (q queue) CountConsumers() int { return len(q.csms) }

// ConsumerAt gets at index
func (q queue) ConsumerAt(i int) Consumer { return q.csms[i] }

// WithName sets name
func (q queue) WithName(name string) Queue {
	q.name = name
	return q
}

// AddTopic adds a topic
func (q queue) AddTopic(tp Topic) Queue {
	q.tps = append(q.tps, tp)
	return q
}

// AddConsumer adds a consumer
func (q queue) AddConsumer(csm Consumer) Queue {
	q.csms = append(q.csms, csm)
	return q
}
