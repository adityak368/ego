module github.com/adityak368/ego/broker

go 1.14

replace github.com/adityak368/ego/broker => ./

require (
	github.com/golang/protobuf v1.4.3
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20191110163157-d32e6e3b99c4 // indirect
)
