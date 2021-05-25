package grpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/adityak368/ego/registry"
	"github.com/adityak368/ego/server"
	"github.com/adityak368/swissknife/logger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// grpcServer creates a new grpc server
type grpcServer struct {
	options     server.Options
	grpcOptions []grpc.ServerOption
	server      *grpc.Server
}

// Name is the server name
func (s *grpcServer) Name() string {
	return s.options.Name
}

// Address is the bind address
func (s *grpcServer) Address() string {
	return s.options.Address
}

// Init initializes the server
func (s *grpcServer) Init(opts server.Options) error {

	s.server = grpc.NewServer(s.grpcOptions...)
	reflection.Register(s.server)
	s.options = opts
	return nil
}

// Options returns the server options
func (s *grpcServer) Options() server.Options {
	return s.options
}

// Run the server
func (s *grpcServer) Run() error {

	if s.server == nil {
		return errors.New("[GRPC-Server]: Cannot Run. Server not Initialized")
	}

	listener, err := net.Listen("tcp", s.Address())
	if err != nil {
		return err
	}

	addr := listener.Addr().(*net.TCPAddr)

	logger.Info().Msgf("[GRPC-Server]: Server listening on %v", addr)

	// add our service details to the registry if present
	if s.options.Registry != nil {
		s.options.Registry.Register(registry.Entry{
			Name:    s.options.Name,
			Address: s.Address(),
			Version: s.options.Version,
		})
		defer s.options.Registry.Deregister(s.options.Name)
		err := s.options.Registry.Watch()
		defer s.options.Registry.CancelWatch()
		if err != nil {
			return err
		}
	}

	if e := s.server.Serve(listener); e != nil {
		return err
	}

	return nil
}

// Handle returns the internal server of grpc
func (s *grpcServer) Handle() interface{} {
	return s.server
}

// The service implementation
func (s *grpcServer) String() string {
	return fmt.Sprintf("[GRPC-Server]: GRPC-Server '%s' Running on %s Version: %s", s.options.Name, s.Address(), s.options.Version)
}

// New returns a new grpc server
func New(grpcOptions ...grpc.ServerOption) server.Server {
	return &grpcServer{
		grpcOptions: grpcOptions,
	}
}
