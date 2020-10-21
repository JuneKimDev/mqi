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
    WithExchange(mqi.NewExchange("test.exchange").
      AddQueue(mqi.NewQueue("test.q1").
        AddTopic(mqi.NewTopic("test.topic1")).
        AddConsumer(mqi.NewConsumer("test.consumer1").WithFunc(mockConsumerFunc))).
      AddQueue(mqi.NewQueue("test.q2").
        AddTopic(mqi.NewTopic("test.topic2")).
        AddConsumer(mqi.NewConsumer("test.consumer2").WithFunc(mockConsumerFunc))))
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

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/JuneKimDev/landing/tags).

## Authors

- **June Kim** - _Initial work_ - [Github](https://github.com/JuneKimDev)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
