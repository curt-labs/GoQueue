package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/curt-labs/GoQueue/helpers/rabbitmq"
	"github.com/curt-labs/GoQueue/helpers/tracker"
	"github.com/streadway/amqp"
)

var (
	configFile            = flag.String("config-file", "", "consumer configuration file")
	numConcurrentHandlers = flag.Int("concurrent", 3, "number of concurrent handlers")
)

type ConsumerHandler struct {
	AppName      string
	GATrackingID string
}

func main() {
	flag.Parse()

	var consumers []*rabbitmq.Consumer

	configs, err := rabbitmq.LoadConsumersConfig(*configFile)
	if err != nil {
		log.Println(err)
		return
	}

	for indx, config := range configs {
		consumer, err := rabbitmq.NewConsumer(
			fmt.Sprintf("consumer%d", indx+1),
			config.QueueName,
			rabbitmq.Exchange{
				Name:       config.ExchangeName,
				RoutingKey: config.RoutingKey,
			}, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		if consumer != nil {
			handler := &ConsumerHandler{
				AppName:      config.RoutingKey,
				GATrackingID: config.GATrackingID,
			}
			consumer.AddConcurrentHandlers(handler, *numConcurrentHandlers)
			consumers = append(consumers, consumer)
		}
	}

	if len(consumers) > 0 {
		for {
			for _, c := range consumers {
				select {
				case <-c.DoneChan:

				}
			}
		}
	}
}

func (h *ConsumerHandler) HandleMessage(msg *amqp.Delivery) error {
	if msg != nil {
		log.Printf("Got message (%s): %s\n", h.AppName, string(msg.Body))

		var event tracker.Event

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			msg.Nack(false, true)
			return err
		}

		//This is kind of temporary...converts segmentio Track object
		//to our tracker event object
		if event.Action == "" {
			var t tracker.Track
			if err := json.Unmarshal(msg.Body, &t); err != nil {
				msg.Nack(false, true)
				return err
			}
			event.Type = tracker.APIEvent
			event.Action = t.Event
			event.Category = "All"
		}

		switch event.Type {
		case tracker.APIEvent:
			event.TrackingID = h.GATrackingID
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
