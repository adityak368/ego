package rabbitmq

import (
	"fmt"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
)

type rabbitmqSubscriber struct {
	topic        string
	subscription *rabbitmq.Channel
	noWait       bool
}

// Topic returns the subscribed topic
func (s *rabbitmqSubscriber) Topic() string {
	return s.topic
}

// Unsubscribe unsibscribes to the topic
func (s *rabbitmqSubscriber) Unsubscribe() error {
	if s.subscription == nil {
		return fmt.Errorf("[RABBITMQ]: Cannot unsubscribe from %s", s.topic)
	}
	return s.subscription.Cancel("", s.noWait)
}
