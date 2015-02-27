package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/curt-labs/GoQueue/helpers/rabbitmq"
	"github.com/curt-labs/GoQueue/helpers/tracker"
	"github.com/streadway/amqp"
)

type ConsumerHandler struct{}

func (h *ConsumerHandler) HandleMessage(msg *amqp.Delivery) error {
	if msg != nil {
		log.Printf("Got message: %s\n", string(msg.Body))

		var event tracker.Event

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			msg.Nack(false, true)
			return err
		}

		switch event.Type {
		case tracker.APIEvent:
			if err := event.SendToGoogleAnalytics(); err != nil {
				msg.Nack(false, true)
				return err
			}
		case tracker.Transaction:
			msg.Nack(false, true)
			return errors.New("Transaction event not implemented!")
		case tracker.EventNotSet:
			msg.Nack(false, true)
			return errors.New("Event type not set")
		default:
			msg.Nack(false, true)
			return errors.New("Unknown event type")
		}

		//acknowledge!
		msg.Ack(false)
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
