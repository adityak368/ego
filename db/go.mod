module github.com/adityak368/ego/db

go 1.14

replace github.com/adityak368/ego/db => ./

require (
	github.com/adityak368/swissknife/logger/v2 v2.0.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	go.mongodb.org/mongo-driver v1.4.2
	google.golang.org/protobuf v1.22.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)
