package server

// Server is an interface for a rpc server
type Server interface {
	// Name is the server name
	Name() string
	// Address is the bind address
	Address() string
	// Init initializes the server
	Init(opts Options) error
	// Options returns the server options
	Options() Options
	// Handle returns the internal server of the implementation
	Handle() interface{}
	// Run the server
	Run() error
	// String returns the description of the server
	String() string
}
