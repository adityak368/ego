module github.com/adityak368/ego/broker

go 1.14

replace github.com/adityak368/ego/broker => ./

replace github.com/adityak368/ego/proto => ../proto

require (
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/isayme/go-amqp-reconnect v0.0.0-20210303120416-fc811b0bcda2
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/wagslane/go-rabbitmq v0.5.1
	google.golang.org/protobuf v1.23.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)
