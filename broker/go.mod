module github.com/adityak368/ego/broker

go 1.18

replace github.com/adityak368/ego/broker => ./

replace github.com/adityak368/ego/proto => ../proto

require (
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	github.com/nats-io/nats.go v1.10.0
	github.com/pkg/errors v0.9.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/wagslane/go-rabbitmq v0.5.1
	google.golang.org/protobuf v1.23.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/jwt v0.3.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nkeys v0.1.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.22.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
