package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adityak368/ego/client"
	"github.com/adityak368/swissknife/logger"
	"google.golang.org/grpc"
)

// GRPC creates a new grpc client
type grpcClient struct {
	options     client.Options
	grpcOptions []grpc.DialOption
	conn        *grpc.ClientConn
}

// Name returns the service name the client connects to
func (g *grpcClient) Name() string {
	return g.options.Name
}

// Address Returns the Target address
func (g *grpcClient) Address() string {
	return g.options.Target
}

// Init initializes the rpc client
func (g *grpcClient) Init(opts client.Options) error {
	g.options = opts
	return nil
}

// Options returns the client options
func (g *grpcClient) Options() client.Options {
	return g.options
}

// Connect connects the client to the rpc server
func (g *grpcClient) Connect() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Infof("[GRPC-Client]: Connecting to %s on %s", g.options.Name, g.options.Target)
	conn, err := grpc.DialContext(
		ctx,
		g.options.Target,
		g.grpcOptions...,
	)
	if err != nil {
		return err
	}

	logger.Infof("[GRPC-Client]: Connected to %s on %s", g.options.Name, g.options.Target)
	g.conn = conn
	return nil
}

// Disconnect disconnects the client
func (g *grpcClient) Disconnect() error {
	if g.conn == nil {
		return errors.New("[GRPC-Client]: Cannot Disconnect. Client not Initialized")
	}
	logger.Infof("[GRPC-Client]: Disconnecting %s from %s", g.options.Name, g.options.Target)
	return g.conn.Close()
}

// String returns the description of the client
func (g *grpcClient) String() string {
	return fmt.Sprintf("[GRPC-Client]: Client %s with target %s", g.options.Name, g.options.Target)
}

// Handle returns the raw connection handle to the rpc server
func (g *grpcClient) Handle() interface{} {
	return g.conn
}

// New creates a new grpc client
func New(grpcOptions ...grpc.DialOption) client.Client {
	return &grpcClient{
		grpcOptions: grpcOptions,
	}
}
