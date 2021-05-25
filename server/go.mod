module github.com/adityak368/ego/server

go 1.14

replace github.com/adityak368/ego/server => ./

require (
	github.com/adityak368/ego/registry v0.0.0-20201103215517-ac54f96660d3
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	google.golang.org/grpc v1.33.0
)
