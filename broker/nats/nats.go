package nats

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/adityak368/ego/broker"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

// Nats is the NATS implementation of the broker
type natsBroker struct {
	options         broker.Options
	connection      *nats.Conn
	subscriptionMap map[string]*natsSubscriber
}

// Address Returns the broker bind interface
func (n *natsBroker) Address() string {
	return n.options.Address
}

// Init initialises the broker
func (n *natsBroker) Init(opts broker.Options) error {
	n.options = opts
	return nil
}

// Options returns the broker options
func (n *natsBroker) Options() broker.Options {
	return n.options
}

// String returns the description of the broker
func (n *natsBroker) String() string {
	return fmt.Sprintf("[NATS]: Connected to NATS on %s", n.Address())
}

// Connect connects to the broker
func (n *natsBroker) Connect() error {
	conn, err := nats.Connect(n.Address())
	if err != nil {
		return err
	}
	log.Printf("[NATS]: Connected to %s", n.Address())
	n.connection = conn
	return nil
}

// Disconnect disconnects from the broker
func (n *natsBroker) Disconnect() error {

	if n.connection == nil {
		return errors.New("[NATS]: Cannot Disconnect. Not connected to broker")
	}

	n.connection.Close()
	log.Printf("[NATS]: Disconnected from %s", n.Address())
	return nil
}

// Handle returns the raw connection handle to the db
func (n *natsBroker) Handle() interface{} {
	return n.connection
}

// Publish publishes a message to the topic
func (n *natsBroker) Publish(topic string, m proto.Message) error {

	if n.connection == nil {
		return errors.New("[NATS]: Cannot Publish. Not connected to broker")
	}

	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	n.connection.Publish(topic, data)
	return nil
}

// PublishRaw publishes raw data to the topic
func (n *natsBroker) PublishRaw(topic string, m []byte) error {

	if n.connection == nil {
		return errors.New("[NATS]: Cannot PublishRaw. Not connected to broker")
	}

	n.connection.Publish(topic, m)
	return nil
}

// Subscribe subscribes a handler to the topic
func (n *natsBroker) Subscribe(topic string, h interface{}) (broker.Subscriber, error) {

	if n.connection == nil {
		return nil, errors.New("[NATS]: Cannot Subscribe. Not connected to broker")
	}

	typ := reflect.TypeOf(h)
	if typ.Kind() != reflect.Func {
		return nil, errors.New("[NATS]: Need a function as a callback")
	}

	if typ.NumIn() != 1 {
		return nil, errors.New("[NATS]: Function should take a single input parameter which is the message")
	}

	msgType := typ.In(0)
	if msgType.Kind() != reflect.Ptr {
		return nil, errors.New("[NATS]: Message should be a pointer")
	}

	cb := reflect.ValueOf(h)

	subscription, err := n.connection.Subscribe(topic, func(m *nats.Msg) {
		msg := reflect.New(msgType.Elem())
		protoMsg, ok := msg.Interface().(proto.Message)
		if !ok {
			log.Println("[NATS]: Message does not implement protobuf message")
			return
		}

		err := proto.Unmarshal(m.Data, protoMsg)
		if err != nil {
			log.Println("[NATS]: Could not decode message")
			return
		}
		cb.Call([]reflect.Value{msg})
	})
	if err != nil {
		return nil, err
	}

	subscriber := &natsSubscriber{
		topic:        topic,
		subscription: subscription,
	}

	n.subscriptionMap[topic] = subscriber
	log.Printf("[NATS]: Subscribed to topic '%s'", topic)
	return subscriber, nil
}

// New returns a new natsBroker broker
func New() broker.Broker {
	return &natsBroker{
		subscriptionMap: make(map[string]*natsSubscriber),
	}
}
