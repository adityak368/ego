package nats

import (
	"testing"
	"time"

	"github.com/adityak368/ego/broker"
	"github.com/stretchr/testify/require"
)

const timeout = 5 * time.Second

func TestNats(t *testing.T) {

	r := require.New(t)

	c := make(chan bool, 1)

	// User is a protobuf message
	OnUserCreatedRaw := func(msg []byte) error {
		// Handle new user creation
		r.Equal(msg, []byte("Test"), "Wrong data received")
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

	subscription, err := bkr.SubscribeRaw("user.UserCreated", OnUserCreatedRaw)
	r.Nil(err)

	r.NotNil(subscription)

	r.Equal(subscription.Topic(), "user.UserCreated", "user.UserCreated subscription error")

	// Publish the protobuf message to the broker
	bkr.PublishRaw("user.UserCreated", []byte("Test"))

	select {
	case <-c:
	case <-time.After(timeout):
		t.Error("Timed out waiting for message from broker")
	}

}
