module github.com/adityak368/ego/registry

go 1.18

replace github.com/adityak368/ego/registry => ./

require (
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	github.com/grandcat/zeroconf v1.0.0
	github.com/pkg/errors v0.9.1
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/miekg/dns v1.1.27 // indirect
	github.com/rs/zerolog v1.22.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
)
