package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type natsSubscriber struct {
	topic        string
	subscription *nats.Subscription
}

// Topic returns the subscribed topic
func (s *natsSubscriber) Topic() string {
	return s.topic
}

// Unsubscribe unsibscribes to the topic
func (s *natsSubscriber) Unsubscribe() error {
	if s.subscription == nil {
		return fmt.Errorf("[NATS]: Cannot unsubscribe from %s", s.topic)
	}
	return s.subscription.Unsubscribe()
}
