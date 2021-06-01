package rabbitmq

type rabbitmqSubscriber struct {
	topic string
}

// Topic returns the subscribed topic
func (s *rabbitmqSubscriber) Topic() string {
	return s.topic
}

// Unsubscribe unsibscribes to the topic
func (s *rabbitmqSubscriber) Unsubscribe() error {
	return nil
}
