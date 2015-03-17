package handlers

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-nsq"
	"net/http"
	"net/url"
)

type AnalyticsHandler struct {
	Category     string
	TrackingCode string
}

func (c *AnalyticsHandler) HandleMessage(message *nsq.Message) error {
	var event Event

	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}

	if event.Action == "" {
		var t Track
		if err := json.Unmarshal(message.Body, &t); err != nil {
			return err
		}

		event.Type = APIEvent
		event.UserId = t.UserId
		event.Action = t.Event
		event.Category = c.Category
		event.TrackingID = c.TrackingCode
	}

	switch event.Type {
	case APIEvent:
		if err := event.SendToGoogleAnalytics(); err != nil {
			return err
		}
	default:
		return errors.New("Unknown event type")
	}

	message.Finish()
	return nil
}

//these are structs copied from SegmentIO
type Message struct {
	Type      string `json:"type,omitempty"`
	MessageId string `json:"messageId,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	SentAt    string `json:"sentAt,omitempty"`
}

type Track struct {
	Context     map[string]interface{} `json:"context,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	AnonymousId string                 `json:"anonymousId,omitempty"`
	UserId      string                 `json:"userId,omitempty"`
	Event       string                 `json:"event"`
	Message
}

//events...so eventful
type EventType int

const (
	EventNotSet EventType = iota
	APIEvent
	Transaction
)

type Event struct {
	Action     string    `json:"action"`
	Category   string    `json:"category"`
	Label      string    `json:"label"`
	TrackingID string    `json:"-"`
	Type       EventType `json:"eventType"`
	UserId     string    `json:"-"`

	Context    map[string]interface{} `json:"context,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func (ev *Event) SendToGoogleAnalytics() error {
	if ev.TrackingID == "" {
		return errors.New("Missing tracking ID")
	}
	vals := make(url.Values, 0)
	vals.Add("v", "1")             //protocol version - required
	vals.Add("tid", ev.TrackingID) //tracking id
	vals.Add("cid", ev.UserId)     //client id - required
	vals.Add("t", "event")         //even type
	vals.Add("ec", ev.Category)    //event category
	vals.Add("ea", ev.Action)      //event action
	vals.Add("el", ev.Label)       //event label
	vals.Add("ev", "200")          //event value

	//more params/detailed information can be found here:
	//https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters

	resp, err := http.Get("http://www.google-analytics.com/collect?" + vals.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
