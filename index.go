package main

import (
	"log"

	"github.com/curt-labs/GoQueue/helpers/rabbitmq"
	"github.com/streadway/amqp"
)

type ConsumerHandler struct {
}

func (h *ConsumerHandler) HandleMessage(msg *amqp.Delivery) error {
	if msg != nil {
		log.Printf("Got message: %s\n", string(msg.Body))
	}
	return nil
}

func main() {
	handler := &ConsumerHandler{}

	exchange := rabbitmq.Exchange{
		Name:       "exchange",
		RoutingKey: "GoAPI",
	}

	consumer, err := rabbitmq.NewConsumer("simple-consumer", "test-queue", exchange, nil)
	if err != nil {
		log.Println(err)
		return
	}

	consumer.AddHandler(handler)

	for {
		select {
		case <-consumer.DoneChan:

		}
	}

	consumer.Close()
}
