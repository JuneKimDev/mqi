package mqi

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func (ch channel) setup() Channel {
	log.Println("Setting up RabbitMQ connection...")

	ch.declareExchange()

	for i := 0; i < ch.Exchange().CountQueues(); i++ {
		q := ch.Exchange().QueueAt(i)
		ch.declareQueue(q)

		for j := 0; j < q.CountTopics(); j++ {
			tp := q.TopicAt(j)
			ch.bindQueueWith(q, tp)
		}

		for j := 0; j < q.CountConsumers(); j++ {
			csm := q.ConsumerAt(j)
			ch.bindConsumerWith(q, csm)
		}
	}

	log.Println("RabbitMQ connection setup is Done")
	return ch
}

func (ch channel) declareExchange() {
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
		os.Exit(1) // Exit b/c failed to setup
	}
	log.Printf("Declared an exchange [%s] successfully\n", exName)
}

func (ch channel) declareQueue(q Queue) {
	qName := q.Name()
	_, err := ch.Sub().QueueDeclare(
		qName, // queue name
		true,  // durable
		false, // auto-delete when there are no remaining consumers
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue [%s]: %v", qName, err) // Exit b/c failed to setup
	}
	log.Printf("Declared a queue [%s] successfully\n", qName)
}

func (ch channel) bindQueueWith(q Queue, tp Topic) {
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
		log.Fatalf("Failed to bind a queue [%s:%s] with [%s]: %v", qName, tpName, exName, err) // Exit b/c failed to setup
	}
	log.Printf("Binded a queue [%s:%s] with [%s] successfully\n", qName, tpName, exName)
}

// runs in goroutine
func (ch channel) bindConsumerWith(q Queue, csm Consumer) {
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
		log.Fatalf("Failed to consume from [%s]: %v", qName, err) // Exit b/c failed to setup
	}
	log.Printf("Consuming a queue [%s]...\n", qName)

	go func() {
		for {
			select {
			case <-ch.KillChan():
				log.Println("Shutting down a consumer")
				return
			case msg := <-msgs:
				log.Printf("MSG[%s]: %s\n", msg.RoutingKey, msg.Body)

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

// sendAck sends ACK to RabbitMQ
func sendAck(msg amqp.Delivery) {
	err := msg.Ack(false)
	if err != nil {
		log.Printf("Failed to send ACK to %d:%v", msg.DeliveryTag, err)
	}
}

// sendNack sends NACK to RabbitMQ
func sendNack(msg amqp.Delivery) {
	err := msg.Nack(false, true)
	if err != nil {
		log.Printf("Failed to send NACK to %d:%v", msg.DeliveryTag, err)
	}
}
