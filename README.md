# EGO Framework

### EGO Framework is a minimalistic opinionated framework for developing microservices in GO.

Although being an opinionated framework, it is still highly extensible and allows you to plug in your own implementations for various components. The framework is heavily inspired by go-micro.
The framework is split into modules and does not bloat the codebase so that you import only the minimal required components for your microservice.

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/adityak368/ego) [![Go Report Card](https://goreportcard.com/badge/github.com/adityak368/ego)](https://goreportcard.com/report/github.com/adityak368/ego)

##### It comes with a project generator to get you started very quickly. It is a Yeoman Generator and you can find it [here](https://www.npmjs.com/package/generator-go-ego)

The Generator uses GRPC + Protobuf + NATS + Kubernetes + Skaffold + Echo + Nginx and supports

- Automatic containerization of services and deployment using kubernetes
- Auto reload on code change
- Clearly defined interfaces for services using protobuf

### Broker

- Defines the broker interface.
- NATS is supported and used by default
- Broker uses protobuf message encoding

```go

    import (
        "log"
        "github.com/adityak368/ego/broker"
        "github.com/adityak368/ego/broker/nats"
        // Replace this with your own protobuf message
        "email/proto/email"
    )

    bkr := nats.New()
    bkr.Init(broker.Options{
        Name: "Nats",
		Address: "localhost:4222",
    })

    err := bkr.Connect()
    if err != nil {
        log.Fatal(err)
    }

    // SendEmailRequest is a protobuf message
    func OnEmail(msg *email.SendEmailRequest) error {
        // Handle new email request
        return nil
    }

    func OnUserCreatedRaw(msg []byte) error {
        // Handle new user creation
        return nil
    }

    emailsubscription, err := bkr.Subscribe("email.SendEmail", OnEmail)
	if err != nil {
        log.Fatal(err)
    }

    usersubscription, err := bkr.SubscribeRaw("user.UserCreated", OnUserCreatedRaw)
    if err != nil {
        log.Fatal(err)
    }

    // Publish the protobuf message to the broker
    bkr.Publish("email.SendEmail", &email.SendEmailRequest{
        Subject: "abcd@example.com",
    })

    // Publish raw message to the broker
    bkr.PublishRaw("user.UserCreated", []byte("Data"))

```

```
syntax = "proto3";

package email;

message SendEmailRequest {
	string To = 1;
	string Subject = 2;
	string Body = 3;
}
```

### Registry

- Defines the Registry interface. Default implementation is mdns. Useful only during local development as most cloud providers block mdns.

```go

    import (
        "log"
        "github.com/adityak368/ego/registry"
        "github.com/adityak368/ego/registry/mdns"
    )

    service := "MyApp"
    serviceName := "MyAwesomeMicroService"
    domain := "local"

    reg := mdns.New(service, domain)
    reg.Init(registry.Options{})

    reg.Register(registry.Entry{
        Name:   serviceName, // Name of the service to register
        Address: "localhost:1212", // Address of the service to register
        Version: "1.0.0",  // Version of the service to register
    })
    err := reg.Watch()
    if err != nil {
        return err
    }
    defer reg.Deregister(serviceName)
    defer reg.CancelWatch()

```

### Server

- Defines the Server interface which serves RPC requests from other microservices
- GRPC is supported and used by default

```go

    import (
        "log"
        // replace with your own service protobuf definition
        "myawesomeapp/proto/myawesomeapp"

        "github.com/adityak368/ego/registry/mdns"
        "github.com/adityak368/ego/server"
        grpcServer "github.com/adityak368/ego/server/grpc"
    )

    func CreateUser(ctx context.Context, req *myawesomeapp.CreateUserRequest) (*myawesomeapp.CreateUserResponse, error) {
        return &myawesomeapp.CreateUserResponse{
            Success: true,
        }, nil
    }

    // Create and start a new grpc server
    srv := grpcServer.New()
    srv.Init(server.Options{
        Name:     "MyAwesomeApp",
        Address:  "localhost:4003",
        Registry: mdns.New("ego", "local"), // Set registry to mdns for automatic service discovery
        Version:  "1.0.0",
    })

    grpcHandle := srv.Handle().(*grpc.Server)

    // Register the protobuf service with the grpc server
    myawesomeapp.RegisterMyAwesomeAppService(grpcHandle, &myawesomeapp.MyAwesomeAppService{
        CreateUser: CreateUser,
    })

```

```
syntax = "proto3";

package myawesomeapp;

service MyAwesomeApp {
	rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
}

message CreateUserRequest {
	string Name = 1;
	string Details = 2;
}

message CreateUserResponse {
	bool success = 1;
}

message UserCreated {
	string Name = 1;
	string Details = 2;
}
```

### Client

- Defines the Client interface which makes a RPC
- GRPC is supported and used by default

```go

    import (
        "log"
        // replace with your own service protobuf definition
        "anotherservice/proto/anotherservice"

        "github.com/adityak368/ego/client"
        grpcClient "github.com/adityak368/ego/client/grpc"
        "github.com/adityak368/mdnsresolver"
        "google.golang.org/grpc"
    )

    anotherServiceClient := grpcClient.New(
        grpc.WithInsecure(),
        grpc.WithResolvers(mdnsresolver.NewBuilder()),
        grpc.WithBalancerName("round_robin"),
    )

    // Initialize the client
    anotherServiceClient.Init(client.Options{
        Name:   "AnotherServiceClient", // Initialize the client by giving it a name
        Target: "mdns://ego/AnotherService.local", // Set service discovery mechanism. MDNS is used here
    })

    // Connect the client to server
    err := anotherServiceClient.Connect()
    if err != nil {
        log.Fatal(err)
    }
    conn := anotherServiceClient.Handle().(*grpc.ClientConn)

    anotherServiceClient = anotherservice.NewAnotherServiceClient(conn)

    // Do RPC Calls
    // anotherServiceClient.SendEmail(...)

```

```
syntax = "proto3";

package anotherservice;

service AnotherService {
	rpc SendEmail(SendEmailRequest) returns (SendEmailResponse) {}
}

message SendEmailRequest {
	string To = 1;
	string Subject = 2;
	string Body = 3;
}

message SendEmailResponse {
	bool success = 1;
}
```

### DB

- Defines the Database and Model interface for connecting to the database
- Contains implementations for MongoDB and Redis

```go

    import (
        "log"
        "github.com/adityak368/ego/db"
        "github.com/adityak368/ego/db/mongodb"
    )

    // MongoDB exports the mongodb handle

    MongoDB := mongodb.New()
    MongoDB.Init(db.Options{
        Name:     "Mongodb",
        Address:  "localhost:27017",
        Database: "MyDatabase",
    })
    err := MongoDB.Connect()
    if err != nil {
        log.Fatal(err)
    }

    userModel := UserModel()
    userModel.CreateIndexes(MongoDB)
    userModel.PrintIndexes(MongoDB)

    defer MongoDB.Disconnect()

```

The User Model

```go

    import (
        "context"
        "time"
        "log"

        "github.com/adityak368/ego/db"
        "github.com/adityak368/ego/db/mongodb"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
    )

    // User defines the user model
    type User struct {
        ID      primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
        Name    string             `bson:"name,omitempty" json:"name" validate:"required"`
        Details string             `bson:"details,omitempty" json:"details" validate:"required"`
    }

    // CreateIndexes creates the indexes for a model
    func (u *User) CreateIndexes(db db.Database) error {
        c := db.Handle().(*mongo.Database).Collection(u.String())
        opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

        keys := bson.D{{"name", 1}}
        index := mongo.IndexModel{}
        index.Keys = keys
        index.Options = &options.IndexOptions{Unique: &[]bool{true}[0]}

        c.Indexes().CreateOne(context.Background(), index, opts)
        return nil
    }

    // PrintIndexes prints all the indexes for the model
    func (u *User) PrintIndexes(db db.Database) {
        switch db.(type) {
        case *mongodb.DB:
            d := db.(*mongodb.DB)
            d.PrintIndexes(u.String())
        default:
            log.Println("Could not load indexes")
        }
    }

    // String returns the string representation of the model
    func (u *User) String() string {
        return "User"
    }

    // UserModel returns user as a db model
    func UserModel() db.Model {
        return &User{}
    }

```
