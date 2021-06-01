package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type rabbitmqSubscriber struct {
	topic        string
	subscription *amqp.Channel
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
