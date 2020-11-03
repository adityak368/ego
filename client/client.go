package client

// Client is an interface used for rpc calls
type Client interface {
	// Name returns the service name the client connects to
	Name() string
	// Init initializes the rpc client
	Init(opts Options) error
	// Options Returns the client options
	Options() Options
	// Address Returns the Target address
	Address() string
	// Connect connects the client to the rpc server
	Connect() error
	// Disconnect disconnects the client
	Disconnect() error
	// Handle returns the raw connection handle to the rpc server
	Handle() interface{}
	// String returns the description of the client
	String() string
}
