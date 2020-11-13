# MQI

Message Queue Interface(MQI) using AMQP

[![PkgGoDev](https://pkg.go.dev/badge/github.com/JuneKimDev/mqi)](https://pkg.go.dev/github.com/JuneKimDev/mqi)
[![Go Report Card](https://goreportcard.com/badge/github.com/JuneKimDev/mqi)](https://goreportcard.com/report/github.com/JuneKimDev/mqi)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/JuneKimDev/mqi)
![GitHub](https://img.shields.io/github/license/JuneKimDev/mqi)

---

## Getting Started

### Prerequisite

You need an active running instance of RabbitMQ server somewhere you can access

### Installing

go get it (pun intended :smile_cat:)

```shell
go get github.com/JuneKimDev/mqi
```

## Usage

```golang
package main

import (
  "log"
  "sync"

  "github.com/JuneKimDev/mqi"
  "github.com/streadway/amqp"
)

// Consumer worker function
func mockConsumerFunc(msg amqp.Delivery) error {
  log.Println("Inside of mockConsumerFunc | Test consumer function")
  log.Println(msg.Body)
  err := mqi.Publish("test.exchangeAnother", "test.topic3", amqp.Publishing{
    ContentType: "text/plain",
    Body:        []byte("Great job")})
  if err != nil {
    return err
  }
  return nil
}

// RabbitMQ setup
func getMockChannelWithConsumer() mqi.Channel {
  return mqi.GetChannel().
    WithBroadcast(NewBroadcast("test.broadcast").
      AddQueue(NewBroadcastQueue().
        AddTopic(NewTopic("test.bctp")).
        AddConsumer(NewConsumer("test.bccsm").WithFunc(mockFunc)))).
    WithExchange(NewExchange("test.exchange").
      AddQueue(NewQueue("test.q1").
        AddTopic(NewTopic("test.topic1")).
        AddConsumer(NewConsumer("test.consumer1").WithFunc(mockFunc))).
      AddQueue(NewQueue("test.q2").
        AddTopic(NewTopic("test.topic2")).
        AddConsumer(NewConsumer("test.consumer2").WithFunc(mockFunc))))
}

func main() {
  var forever sync.WaitGroup
  forever.Add(1)
  // Connect to RabbitMQ
  getMockChannelWithConsumer().Start()
  log.Println("The app is running in forever loop")
  forever.Wait() // Prevents service from exiting
}
```
