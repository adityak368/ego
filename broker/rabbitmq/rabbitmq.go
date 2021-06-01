package rabbitmq

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/adityak368/ego/broker"
	"github.com/adityak368/swissknife/logger/v2"
	"github.com/streadway/amqp"
)

// RabbitMq is the RABBITMQ implementation of the broker
type rabbitmqBroker struct {
	options         broker.Options
	connection      *amqp.Connection
	subscriptionMap map[string]*rabbitmqSubscriber
	queueMap        map[string]*amqp.Channel
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

	conn, err := amqp.Dial(n.Address())
	if err != nil {
		return err
	}

	logger.Info().Msgf("[RABBITMQ]: Connected to %s", n.Address())
	n.connection = conn

	return nil
}

// Disconnect disconnects from the broker
func (n *rabbitmqBroker) Disconnect() error {

	if n.connection == nil {
		return errors.New("[RABBITMQ]: Cannot Disconnect. Not connected to broker")
	}

	for _, v := range n.queueMap {
		err := v.Close()
		if err != nil {
			return err
		}
	}

	err := n.connection.Close()
	if err != nil {
		return err
	}

	logger.Info().Msgf("[RABBITMQ]: Disconnected from %s", n.Address())
	return nil
}

// Handle returns the raw connection handle to the db
func (n *rabbitmqBroker) Handle() interface{} {
	return n.connection
}

// Publish publishes a message to the topic
func (n *rabbitmqBroker) Publish(topic string, m proto.Message) error {

	if n.connection == nil {
		return errors.New("[RABBITMQ]: Cannot Publish. Not connected to broker")
	}

	var ch *amqp.Channel
	ch, ok := n.queueMap[topic]
	if !ok {
		newChannel, err := n.connection.Channel()
		if err != nil {
			return err
		}
		_, err = ch.QueueDeclare(
			topic,
			n.config.Durable,
			n.config.DeleteWhenUnused,
			n.config.Exclusive,
			n.config.NoWait,
			n.config.Arguments,
		)
		if err != nil {
			return err
		}
		n.queueMap[topic] = newChannel
		ch = newChannel
	}

	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	err = ch.Publish(
		n.config.Exchange,
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        data,
		},
	)

	return err
}

// PublishRaw publishes raw data to the topic
func (n *rabbitmqBroker) PublishRaw(topic string, m []byte) error {

	if n.connection == nil {
		return errors.New("[RABBITMQ]: Cannot Publish. Not connected to broker")
	}

	var ch *amqp.Channel
	ch, ok := n.queueMap[topic]
	if !ok {
		newChannel, err := n.connection.Channel()
		if err != nil {
			return err
		}
		_, err = ch.QueueDeclare(
			topic,
			n.config.Durable,
			n.config.DeleteWhenUnused,
			n.config.Exclusive,
			n.config.NoWait,
			n.config.Arguments,
		)
		if err != nil {
			return err
		}
		n.queueMap[topic] = newChannel
		ch = newChannel
	}

	err := ch.Publish(
		n.config.Exchange,
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        m,
		},
	)

	return err
}

// Subscribe subscribes a handler to the topic
func (n *rabbitmqBroker) Subscribe(topic string, h interface{}) (broker.Subscriber, error) {

	if n.connection == nil {
		return nil, errors.New("[RABBITMQ]: Cannot Subscribe. Not connected to broker")
	}

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

	var ch *amqp.Channel
	ch, ok := n.queueMap[topic]
	if !ok {
		newChannel, err := n.connection.Channel()
		if err != nil {
			return nil, err
		}
		ch = newChannel
		n.queueMap[topic] = newChannel
		_, err = ch.QueueDeclare(
			topic,
			n.config.Durable,
			n.config.DeleteWhenUnused,
			n.config.Exclusive,
			n.config.NoWait,
			n.config.Arguments,
		)
		if err != nil {
			return nil, err
		}
	}

	msgs, err := ch.Consume(
		topic, // queue
		"",    // consumer
		n.config.AutoAck,
		n.config.Exclusive,
		false, // no-local
		n.config.NoWait,
		n.config.Arguments,
	)

	if err != nil {
		return nil, err
	}

	go func() {
		for d := range msgs {
			msg := reflect.New(msgType.Elem())

			protoMsg, ok := msg.Interface().(proto.Message)
			if !ok {
				logger.Warn().Msg("[RABBITMQ]: Message does not implement protobuf message")
				return
			}

			err := proto.Unmarshal(d.Body, protoMsg)
			if err != nil {
				logger.Warn().Msg("[RABBITMQ]: Could not decode message")
				return
			}

			res := cb.Call([]reflect.Value{reflect.ValueOf(context.Background()), msg})

			if len(res) != 1 {
				logger.Warn().Msg("[RABBITMQ]: Invalid return value")
				return
			}

			if v := res[0].Interface(); v != nil {
				err, ok := v.(error)
				if !ok {
					logger.Warn().Msg("[RABBITMQ]: Could not parse error")
				} else {
					logger.Error().Err(err).Msg("")
				}
				continue
			}

			if !n.config.AutoAck {
				d.Ack(false)
			}
		}
	}()

	subscriber := &rabbitmqSubscriber{
		topic:        topic,
		subscription: ch,
		noWait:       n.config.NoWait,
	}

	n.subscriptionMap[topic] = subscriber
	logger.Info().Msgf("[RABBITMQ]: Subscribed to topic '%s'", topic)
	return subscriber, nil
}

// SubscribeRaw subscribes a raw handler to the topic
func (n *rabbitmqBroker) SubscribeRaw(topic string, h func(c context.Context, data []byte) error) (broker.Subscriber, error) {

	if n.connection == nil {
		return nil, errors.New("[RABBITMQ]: Cannot Subscribe. Not connected to broker")
	}

	var ch *amqp.Channel
	ch, ok := n.queueMap[topic]
	if !ok {
		newChannel, err := n.connection.Channel()
		if err != nil {
			return nil, err
		}
		ch = newChannel
		n.queueMap[topic] = newChannel
		_, err = ch.QueueDeclare(
			topic,
			n.config.Durable,
			n.config.DeleteWhenUnused,
			n.config.Exclusive,
			n.config.NoWait,
			n.config.Arguments,
		)
		if err != nil {
			return nil, err
		}
	}

	msgs, err := ch.Consume(
		topic, // queue
		"",    // consumer
		n.config.AutoAck,
		n.config.Exclusive,
		false, // no-local
		n.config.NoWait,
		n.config.Arguments,
	)

	if err != nil {
		return nil, err
	}

	go func() {
		for d := range msgs {
			err := h(context.Background(), d.Body)
			if err != nil {
				logger.Error().Err(err).Msg("")
				continue
			}

			if !n.config.AutoAck {
				d.Ack(false)
			}
		}
	}()

	subscriber := &rabbitmqSubscriber{
		topic:        topic,
		subscription: ch,
		noWait:       n.config.NoWait,
	}

	n.subscriptionMap[topic] = subscriber
	logger.Info().Msgf("[RABBITMQ]: Subscribed to topic '%s'", topic)
	return subscriber, nil
}

// New returns a new rabbitmqBroker broker
func New(config Config) broker.Broker {
	return &rabbitmqBroker{
		subscriptionMap: make(map[string]*rabbitmqSubscriber),
		config:          config,
		queueMap:        make(map[string]*amqp.Channel),
	}
}
