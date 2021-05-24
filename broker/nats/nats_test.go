package nats

import (
	"context"
	"testing"
	"time"

	"github.com/adityak368/ego/broker"
	proto "github.com/adityak368/ego/broker/proto/gen/broker"
	"github.com/stretchr/testify/require"
)

const timeout = 5 * time.Second

func TestNats(t *testing.T) {

	r := require.New(t)

	c := make(chan bool, 1)

	// msg is a raw message
	OnTestMessageRaw := func(msg []byte) error {
		// Handle new user creation
		r.Equal(msg, []byte("Test"), "Wrong data received")
		c <- true
		return nil
	}

	// TestMessage is a protobuf message
	OnTestMessageProto := func(ctx context.Context, msg *proto.TestMessage) error {
		// Handle new user creation
		r.Equal(msg.Data, "Test", "Wrong data received")
		c <- true
		return nil
	}

	bkr := New()
	bkr.Init(broker.Options{
		Name:    "Nats",
		Address: "localhost:4222",
	})

	err := bkr.Connect()
	r.Nil(err)

	subscriptionRaw, err := bkr.SubscribeRaw("test.testMessageRaw", OnTestMessageRaw)
	r.Nil(err)
	r.NotNil(subscriptionRaw)
	subscriptionProto, err := bkr.Subscribe("test.testMessageProto", OnTestMessageProto)
	r.Nil(err)
	r.NotNil(subscriptionProto)

	r.Equal(subscriptionRaw.Topic(), "test.testMessageRaw", "test.testMessageRaw subscription error")
	r.Equal(subscriptionProto.Topic(), "test.testMessageProto", "test.testMessageProto subscription error")

	// Publish the protobuf message to the broker
	bkr.PublishRaw("test.testMessageRaw", []byte("Test"))
	bkr.Publish("test.testMessageProto", &proto.TestMessage{Data: "Test"})

	select {
	case <-c:
	case <-time.After(timeout):
		t.Error("Timed out waiting for message from broker")
	}

}
