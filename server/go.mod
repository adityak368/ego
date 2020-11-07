module github.com/adityak368/ego/server

go 1.14

replace github.com/adityak368/ego/server => ./

require (
	github.com/adityak368/ego/registry v0.0.0-20201103215517-ac54f96660d3
	github.com/adityak368/swissknife/logger v0.0.0-20201107160000-5f5e30188eb2
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/grpc v1.33.0
)
