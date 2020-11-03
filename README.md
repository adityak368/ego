# EGO Framework

### EGO Framework is a minimalistic opinionated framework for developing microservices in GO.

Although being an opinionated framework, it still highly extensible and allows you to plug in your own implementations for various components. The framework is heavily inspired by go-micro.
The framework is split into modules and does not bloat the codebase so that you import only the minimal required components for your microservice.

It comes with a generator project to get you started very quickly.
It uses GRPC + Protobuf + NATS + Kubenetes + Skaffold + Echo

```
go get -u github.com/adityak368/ego/<modulename>@main

```

### Broker

- Defines the broker interface. Default implementation is NATS
- Broker uses protobuf message encoding

```go

    import (
        "github.com/adityak368/ego/broker"
        "github.com/adityak368/ego/broker/nats"
        // Replace with your own protobuf message
        "github.com/adityak368/ego/test/email"
    )

    bkr := nats.New()
    err = bkr.Init(broker.Options{
        Name: "Nats",
        Host: "localhost",
        Port: 4222,
    })
    if err != nil {
        log.Fatal(err)
    }
    err = bkr.Connect()
    if err != nil {
        log.Fatal(err)
    }

    // Publish the message to the broker
    bkr.Publish("sendEmail", &email.SendEmailRequest{
        Subject: "abcd@example.com",
    })

```

### Registry

- Defines the Registry interface. Default implementation is mdns. Useful only during local development as most cloud providers block mdns.

```go

    import (
        "github.com/adityak368/ego/registry"
        "github.com/adityak368/ego/registry/mdns"
    )

    service := "MyApp"
    serviceName := "MyAwesomeMicroService"
    domain := "local"

    reg := mdns.New(service, domain)
    reg.Init(registry.Options{})

    reg.Register(registry.Entry{
        Name:   serviceName,
        Address: "localhost:1212",
        Version: "1.0.0",
    })
    err := reg.Watch()
    if err != nil {
        return err
    }
    defer reg.Deregister(serviceName)
    defer reg.CancelWatch()

```

### Server

- Defines the Server interface which serves RPC from other microservices
- GRPC is supported

```go

    import (
        "sampleservice/config"
        // service protobuf definition
        "sampleservice/proto/sampleservice"

        "github.com/adityak368/ego/server"
        grpcServer "github.com/adityak368/ego/server/grpc"
    )

    func CreateUser(ctx context.Context, req *sampleservice.CreateUserRequest) (*sampleservice.CreateUserResponse, error) {
        return &sampleservice.CreateUserResponse{
            Success: true,
        }, nil
    }

    // Create and start a new grpc server
    srv := grpcServer.New(
    // grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
    // grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(tracer)),
    )
    srv.Init(server.Options{
        Name:     config.AppName,
        Address:  config.AddressMicroservice,
        Registry: mdns.New("ego", "local"),
        Version:  "1.0.0",
    })

    grpcHandle := srv.Handle().(*grpc.Server)
    handler := sampleservice.SampleServiceService{
        CreateUser: CreateUser,
    }

    // Register the protobuf service with the grpc server
    sampleservice.RegisterSampleServiceService(grpcHandle, &handler)

```

### Client

- Defines the Client interface which makes a RPC
- GRPC is supported

```go

    import (
        // the protobuf service interface
        "anotherservice/proto/anotherservice"

        "github.com/adityak368/ego/client"
        grpcClient "github.com/adityak368/ego/client/grpc"
        "google.golang.org/grpc"
    )

    anotherServiceClient := grpcClient.New(
        grpc.WithInsecure(),
        grpc.WithResolvers(mdnsresolver.NewBuilder()),
        grpc.WithBalancerName("round_robin"),
    )

    // Initialize the client
    anotherServiceClient.Init(client.Options{
        Name:   "AnotherService",
        Target: "dns://AnotherService.local",
    })

    // Connect the client to server
    err := anotherServiceClient.Connect()
    if err != nil {
        log.Fatal(err)
    }
    conn := anotherServiceClient.Handle().(*grpc.ClientConn)

    anotherServiceClient = anotherservice.NewAnotherServiceClient(conn)

    // Do RPC Calls
    // anotherServiceClient.SendEmail()

```

### DB

- Defines the Database and Model interface for connecting to the database
- Contains implementations for MongoDB and Redis

```go

    import (
        "sampleservice/config"
        "db/models"

        "github.com/adityak368/ego/db"
        "github.com/adityak368/ego/db/mongodb"
        "github.com/adityak368/swissknife/logger"
    )

    // MongoDB exports the mongodb handle

    MongoDB := mongodb.New()
    MongoDB.Init(db.Options{
        Name:     "Mongodb",
        Address:  config.MongoDBUrl,
        Database: config.DbName,
    })
    err := MongoDB.Connect()
    if err != nil {
        log.Fatal(err)
    }

    userModel := models.UserModel()
    userModel.CreateIndexes(MongoDB)
    userModel.PrintIndexes(MongoDB)

    defer MongoDB.Disconnect()

```

The User Model

```go

    package models

    import (
        "context"
        "sampleservice/config"
        "time"

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
        opts := options.CreateIndexes().SetMaxTime(time.Duration(config.DbTimeout) * time.Second)

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
