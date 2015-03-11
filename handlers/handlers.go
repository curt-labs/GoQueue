package handlers

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-nsq"
	"github.com/curt-labs/GoQueue/tracker"
)

type ConsumerHandler struct {
	Category     string
	TrackingCode string
}

func (c *ConsumerHandler) HandleMessage(message *nsq.Message) error {
	var event tracker.Event

	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}

	if event.Action == "" {
		var t tracker.Track
		if err := json.Unmarshal(message.Body, &t); err != nil {
			return err
		}

		event.Type = tracker.APIEvent
		event.UserId = t.UserId
		event.Action = t.Event
		event.Category = c.Category
		event.TrackingID = c.TrackingCode
	}

	switch event.Type {
	case tracker.APIEvent:
		if err := event.SendToGoogleAnalytics(); err != nil {
			return err
		}
	default:
		return errors.New("Unknown event type")
	}

	message.Finish()
	return nil
}
