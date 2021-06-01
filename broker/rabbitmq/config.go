package rabbitmq

import "github.com/streadway/amqp"

type Config struct {
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	AutoAck          bool
	Exchange         string
	Arguments        amqp.Table
}
