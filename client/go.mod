module github.com/adityak368/ego/client

go 1.18

replace github.com/adityak368/ego/client => ./

require (
	github.com/adityak368/ego/registry v1.0.1
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.33.0
)

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/rs/zerolog v1.22.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
	google.golang.org/protobuf v1.23.0 // indirect
)
