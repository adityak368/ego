package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/adityak368/ego/broker"
	proto "github.com/adityak368/ego/broker/proto/gen/broker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const timeout = 5 * time.Second

func TestNats(t *testing.T) {

	r := require.New(t)

	c := make(chan bool, 1)

	// msg is a raw message
	OnTestMessageRaw := func(ctx context.Context, msg []byte) error {
		// Handle new user creation
		r.Equal(msg, []byte("Test"), "Wrong data received")
		return nil
	}

	// TestMessage is a protobuf message
	OnTestMessageProto := func(ctx context.Context, msg *proto.TestMessage) error {
		// Handle new user creation
		r.Equal(msg.Data, "Test", "Wrong data received")
		return nil
	}

	// TestMessage is a protobuf message
	OnTestMessageProtoWithError := func(ctx context.Context, msg *proto.TestMessage) error {
		// Handle new user creation
		r.Equal(msg.Data, "Test", "Wrong data received")
		timer := time.NewTimer(1 * time.Second)

		go func() {
			<-timer.C
			c <- true
		}()

		return errors.New("Something went wrong")
	}

	bkr := New(Config{
		Durable:          true,
		DeleteWhenUnused: false,
		Exclusive:        false,
		NoWait:           true,
		AutoAck:          false,
	})
	bkr.Init(broker.Options{
		Name:    "RabbitMQ",
		Address: "amqp://localhost:5672",
	})

	err := bkr.Connect()
	r.Nil(err)

	subscriptionRaw, err := bkr.SubscribeRaw("test.testMessageRaw", OnTestMessageRaw)
	r.Nil(err)
	r.NotNil(subscriptionRaw)
	subscriptionProto, err := bkr.Subscribe("test.testMessageProto", OnTestMessageProto)
	r.Nil(err)
	r.NotNil(subscriptionProto)
	subscriptionProtoWithError, err := bkr.Subscribe("test.testMessageProtoWithError", OnTestMessageProtoWithError)
	r.Nil(err)
	r.NotNil(subscriptionProtoWithError)

	r.Equal(subscriptionRaw.Topic(), "test.testMessageRaw", "test.testMessageRaw subscription error")
	r.Equal(subscriptionProto.Topic(), "test.testMessageProto", "test.testMessageProto subscription error")
	r.Equal(subscriptionProtoWithError.Topic(), "test.testMessageProtoWithError", "test.testMessageProtoWithError subscription error")

	// Publish the protobuf message to the broker
	bkr.PublishRaw("test.testMessageRaw", []byte("Test"))
	bkr.Publish("test.testMessageProto", &proto.TestMessage{Data: "Test"})
	bkr.Publish("test.testMessageProtoWithError", &proto.TestMessage{Data: "Test"})

	select {
	case <-c:
	case <-time.After(timeout):
		t.Error("Timed out waiting for message from broker")
	}

}
