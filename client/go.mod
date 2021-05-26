module github.com/adityak368/ego/client

go 1.14

replace github.com/adityak368/ego/client => ./

require (
	github.com/adityak368/ego/registry v1.0.1
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.33.0
)
