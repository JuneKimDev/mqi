package mqi

import (
	"log"

	"github.com/streadway/amqp"
)

func setup() {
	log.Println("Setting up RabbitMQ connection...")

	declareExchange()

	ch := GetChannel()
	for i := 0; i < ch.Exchange().CountQueues(); i++ {
		q := ch.Exchange().QueueAt(i)
		declareQueue(q)

		ch = GetChannel()
		for j := 0; j < q.CountTopics(); j++ {
			tp := q.TopicAt(j)
			bindQueueWith(q, tp)
		}

		for j := 0; j < q.CountConsumers(); j++ {
			csm := q.ConsumerAt(j)
			bindConsumerWith(q, csm)
		}
	}

	log.Println("RabbitMQ connection setup is Done")
}

func declareExchange() {
	ch := GetChannel()
	exName := ch.Exchange().Name()
	err := ch.Sub().ExchangeDeclare(
		exName,  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
		// Exit b/c failed to setup
	}
	log.Printf("Declared an exchange [%s] successfully\n", exName)
}

func declareQueue(q Queue) {
	ch := GetChannel()
	qName := q.Name()
	ref, err := ch.Sub().QueueDeclare(
		qName, // queue name
		true,  // durable; it must follow durability of exchange with which it is binded
		false, // auto-delete when there are no remaining consumers
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue [%s]: %v", qName, err)
		// Exit b/c failed to setup
	}
	ch.UpdateChan() <- ch.WithExchange(ch.Exchange().UpdateQueue(q.WithRef(&ref)))
	log.Printf("Declared a queue [%s] successfully\n", qName)
}

func declareTempQueue(q Queue) Queue {
	ch := GetChannel()
	qName := q.Name()
	// settings for Temporary Queue (auto-delete is on)
	// As soon as the consumer finishes the execution, the queue will be deleted automatically
	ref, err := ch.Sub().QueueDeclare(
		qName, // queue name
		true,  // durable; it must follow durability of exchange with which it is binded
		true,  // auto-delete when there are no remaining consumers
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue [%s]: %v", qName, err)
		// Exit b/c failed to setup
	}
	log.Printf("Declared a queue [%s] successfully\n", qName)
	return q.WithRef(&ref)
}

func bindQueueWith(q Queue, tp Topic) {
	ch := GetChannel()
	exName := ch.Exchange().Name()
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

// runs in goroutine
func bindConsumerWith(q Queue, csm Consumer) {
	ch := GetChannel()
	qName := q.Name()

	msgs, err := ch.Sub().Consume(
		qName,      // queue name
		csm.Name(), // consumer name
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to consume from [%s]: %v", qName, err)
		// Exit b/c failed to setup
	}
	log.Printf("Consuming a queue [%s]...\n", qName)

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

// runs in goroutine; but runs only one time
func bindTempConsumerWith(q Queue, csm Consumer) {
	ch := GetChannel()
	qName := q.Name()

	msgs, err := ch.Sub().Consume(
		qName,      // queue name
		csm.Name(), // consumer name
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to consume from [%s]: %v", qName, err)
		// Exit b/c failed to setup
	}
	log.Printf("Consuming a queue [%s]...\n", qName)

	// Consumer in APIGateway doesn't loop forever
	go func() {
		msg := <-msgs
		if err := csm.Func()(msg); err != nil {
			log.Printf("Error in consumer function: %v\n", err)
			sendNack(msg)
		} else {
			sendAck(msg)
		}
	}()
}

// sendAck sends ACK to RabbitMQ
func sendAck(msg amqp.Delivery) {
	err := msg.Ack(false)
	if err != nil {
		log.Fatalf("Failed to send ACK to %d:%v", msg.DeliveryTag, err)
	}
}

// sendNack sends NACK to RabbitMQ
func sendNack(msg amqp.Delivery) {
	err := msg.Nack(false, true)
	if err != nil {
		log.Fatalf("Failed to send NACK to %d:%v", msg.DeliveryTag, err)
	}
}
