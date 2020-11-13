package mqi

import (
	"log"
)

func setup() {
	log.Println("Setting up RabbitMQ connection...")

	ch := GetChannel()

	// Setup broadcast
	ch.Broadcast().declare()
	bcq := ch.Broadcast().QueueAt(0).declare()
	bcq.TopicAt(0).bindWith(bcq)
	bcq.ConsumerAt(0).bindWith(bcq)

	// Setup normals
	ex := ch.Exchange()
	ex.declare()
	for i := 0; i < ex.CountQueues(); i++ {
		q := ex.QueueAt(i).declare()

		for j := 0; j < q.CountTopics(); j++ {
			q.TopicAt(j).bindWith(q)
		}

		for j := 0; j < q.CountConsumers(); j++ {
			q.ConsumerAt(j).bindWith(q)
		}
	}

	log.Println("RabbitMQ connection setup is Done")
}

// AddTempQueue adds a queue which gets auto-deleted after execution of consumer function
func AddTempQueue(q Queue) Queue {
	log.Println("Adding a queue...")
	dq := q.declare()

	for j := 0; j < dq.CountTopics(); j++ {
		dq.TopicAt(j).bindWith(dq)
	}

	for j := 0; j < dq.CountConsumers(); j++ {
		dq.ConsumerAt(j).bindWith(dq)
	}
	return dq
}
