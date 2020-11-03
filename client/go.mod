module github.com/adityak368/ego/client

go 1.14

replace github.com/adityak368/ego/client => ./

replace github.com/adityak368/ego/registry => ../registry

require (
	github.com/adityak368/ego/registry v0.0.0-00010101000000-000000000000
	github.com/grandcat/zeroconf v1.0.0
	github.com/micro/mdns v0.3.0
	github.com/prometheus/client_golang v1.8.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/grpc v1.33.0
)
