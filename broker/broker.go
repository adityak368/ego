// Package broker is an interface used for asynchronous messaging
package broker

import (
	"github.com/golang/protobuf/proto"
)

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	// Topic returns the subscribed topic
	Topic() string
	// Unsubscribe unsubscribes to the topic
	Unsubscribe() error
}

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	// Init initializes the broker
	Init(opts Options) error
	// Options Returns the broker options
	Options() Options
	// Address Returns the broker bind interface
	Address() string
	// Connect connects to the broker
	Connect() error
	// Disconnect disconnects from the broker
	Disconnect() error
	// Publish publishes a message to the topic
	Publish(topic string, m proto.Message) error
	// Publish publishes raw data to the topic
	PublishRaw(topic string, m []byte) error
	// Subscribe subscribes a handler to the topic
	Subscribe(topic string, h interface{}) (Subscriber, error)
	// SubscribeRaw subscribes a raw handler to the topic
	SubscribeRaw(topic string, h func(data []byte) error) (Subscriber, error)
	// Handle returns the raw connection handle to the broker
	Handle() interface{}
	// String returns the description of the broker
	String() string
}
