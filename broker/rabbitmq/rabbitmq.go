package rabbitmq

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/adityak368/ego/broker"
	"github.com/adityak368/swissknife/logger/v2"
	"github.com/wagslane/go-rabbitmq"
)

// RabbitMq is the RABBITMQ implementation of the broker
type rabbitmqBroker struct {
	options         broker.Options
	subscriptionMap map[string]rabbitmq.Consumer
	publisherMap    map[string]rabbitmq.Publisher
	config          Config
}

// Address Returns the broker bind interface
func (n *rabbitmqBroker) Address() string {
	return n.options.Address
}

// Init initialises the broker
func (n *rabbitmqBroker) Init(opts broker.Options) error {
	n.options = opts
	return nil
}

// Options returns the broker options
func (n *rabbitmqBroker) Options() broker.Options {
	return n.options
}

// String returns the description of the broker
func (n *rabbitmqBroker) String() string {
	return fmt.Sprintf("[RABBITMQ]: Connected to RabbitMQ on %s", n.Address())
}

// Connect connects to the broker
func (n *rabbitmqBroker) Connect() error {
	logger.Info().Msgf("[RABBITMQ]: Connected to %s", n.Address())
	return nil
}

// Disconnect disconnects from the broker
func (n *rabbitmqBroker) Disconnect() error {
	logger.Info().Msgf("[RABBITMQ]: Disconnected from %s", n.Address())
	return nil
}

// Handle returns the raw connection handle to the db
func (n *rabbitmqBroker) Handle() interface{} {
	return nil
}

// Publish publishes a message to the topic
func (n *rabbitmqBroker) Publish(topic string, m proto.Message) error {

	var p rabbitmq.Publisher
	p, ok := n.publisherMap[topic]
	if !ok {
		publisher, _, err := rabbitmq.NewPublisher(
			n.Address(),
			// can pass nothing for no logging
			rabbitmq.WithPublisherOptionsLogging,
		)
		if err != nil {
			return err
		}
		n.publisherMap[topic] = publisher
		p = publisher
	}

	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	err = p.Publish(
		data,
		[]string{topic},
		rabbitmq.WithPublishOptionsContentType("application/octet-stream"),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		func(options *rabbitmq.PublishOptions) {
			options.Exchange = ""
			options.Mandatory = false
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// PublishRaw publishes raw data to the topic
func (n *rabbitmqBroker) PublishRaw(topic string, m []byte) error {

	var p rabbitmq.Publisher
	p, ok := n.publisherMap[topic]
	if !ok {
		publisher, _, err := rabbitmq.NewPublisher(
			n.Address(),
			// can pass nothing for no logging
			rabbitmq.WithPublisherOptionsLogging,
		)
		if err != nil {
			return err
		}
		n.publisherMap[topic] = publisher
		p = publisher
	}

	err := p.Publish(
		m,
		[]string{topic},
		rabbitmq.WithPublishOptionsContentType("application/octet-stream"),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		func(options *rabbitmq.PublishOptions) {
			options.Exchange = ""
			options.Mandatory = false
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Subscribe subscribes a handler to the topic
func (n *rabbitmqBroker) Subscribe(topic string, h interface{}) (broker.Subscriber, error) {

	typ := reflect.TypeOf(h)
	if typ.Kind() != reflect.Func {
		return nil, errors.New("[RABBITMQ]: Need a function as a callback")
	}

	if typ.NumIn() != 2 {
		return nil, errors.New("[RABBITMQ]: Function takes two inputs. 1. context.Context and 2. proto.Message which is the message")
	}

	ctxType := typ.In(0)
	if ctxType.Kind() != reflect.Interface || !ctxType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return nil, errors.New("[RABBITMQ]: First Parameter should be of type context.Context")
	}

	msgType := typ.In(1)
	if msgType.Kind() != reflect.Ptr {
		return nil, errors.New("[RABBITMQ]: Message should be a pointer")
	}

	if typ.NumOut() != 1 {
		return nil, errors.New("[RABBITMQ]: Function should have a single return value")
	}

	errType := typ.Out(0)
	if errType.Kind() != reflect.Interface || !errType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, errors.New("[RABBITMQ]: Function should return error or nil")
	}

	cb := reflect.ValueOf(h)

	var c rabbitmq.Consumer
	c, ok := n.subscriptionMap[topic]
	if !ok {
		consumer, err := rabbitmq.NewConsumer(
			n.Address(),
			rabbitmq.WithConsumerOptionsLogging,
		)
		if err != nil {
			logger.Error().Err(err).Msg("")
		}
		n.subscriptionMap[topic] = consumer
		c = consumer
	}

	err := c.StartConsuming(
		func(d rabbitmq.Delivery) bool {
			msg := reflect.New(msgType.Elem())

			protoMsg, ok := msg.Interface().(proto.Message)
			if !ok {
				logger.Warn().Msg("[RABBITMQ]: Message does not implement protobuf message")
				return false
			}

			err := proto.Unmarshal(d.Body, protoMsg)
			if err != nil {
				logger.Warn().Msg("[RABBITMQ]: Could not decode message")
				return false
			}

			res := cb.Call([]reflect.Value{reflect.ValueOf(context.Background()), msg})

			if len(res) != 1 {
				logger.Warn().Msg("[RABBITMQ]: Invalid return value")
				return false
			}

			if v := res[0].Interface(); v != nil {
				err, ok := v.(error)
				if !ok {
					logger.Warn().Msg("[RABBITMQ]: Could not parse error")
				} else {
					logger.Error().Err(err).Msg("")
				}
				return false
			}

			return true
		},
		topic,
		nil,
		rabbitmq.WithConsumeOptionsConcurrency(1),
		func(options *rabbitmq.ConsumeOptions) {
			options.ConsumerAutoAck = n.config.AutoAck
			options.ConsumerArgs = rabbitmq.Table(n.config.Arguments)
			options.ConsumerNoWait = n.config.NoWait
			options.QueueAutoDelete = n.config.DeleteWhenUnused
			options.QueueExclusive = n.config.Exclusive
			options.QueueDurable = n.config.Durable
		},
	)

	if err != nil {
		return nil, err
	}

	subscriber := &rabbitmqSubscriber{
		topic: topic,
	}

	logger.Info().Msgf("[RABBITMQ]: Subscribed to topic '%s'", topic)
	return subscriber, nil
}

// SubscribeRaw subscribes a raw handler to the topic
func (n *rabbitmqBroker) SubscribeRaw(topic string, h func(c context.Context, data []byte) error) (broker.Subscriber, error) {

	var c rabbitmq.Consumer
	c, ok := n.subscriptionMap[topic]
	if !ok {
		consumer, err := rabbitmq.NewConsumer(
			n.Address(),
			rabbitmq.WithConsumerOptionsLogging,
		)
		if err != nil {
			logger.Error().Err(err).Msg("")
		}
		n.subscriptionMap[topic] = consumer
		c = consumer
	}

	err := c.StartConsuming(
		func(d rabbitmq.Delivery) bool {

			err := h(context.Background(), d.Body)
			if err != nil {
				logger.Error().Err(err).Msg("")
				return false
			}

			return true
		},
		topic,
		nil,
		rabbitmq.WithConsumeOptionsConcurrency(1),
		func(options *rabbitmq.ConsumeOptions) {
			options.ConsumerAutoAck = n.config.AutoAck
			options.ConsumerArgs = rabbitmq.Table(n.config.Arguments)
			options.ConsumerNoWait = n.config.NoWait
			options.QueueAutoDelete = n.config.DeleteWhenUnused
			options.QueueExclusive = n.config.Exclusive
			options.QueueDurable = n.config.Durable
		},
	)

	if err != nil {
		return nil, err
	}

	subscriber := &rabbitmqSubscriber{
		topic: topic,
	}

	logger.Info().Msgf("[RABBITMQ]: Subscribed to topic '%s'", topic)
	return subscriber, nil
}

// New returns a new rabbitmqBroker broker
func New(config Config) broker.Broker {
	return &rabbitmqBroker{
		subscriptionMap: make(map[string]rabbitmq.Consumer),
		publisherMap:    make(map[string]rabbitmq.Publisher),
		config:          config,
	}
}
