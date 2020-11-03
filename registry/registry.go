package registry

// Registry defines the registry interface
type Registry interface {
	// Init initializes the registry
	Init(opts Options) error
	// Options Returns the registry options
	Options() Options
	// Register adds the service to the registry
	Register(entry Entry) error
	// Deregister removes the service to the registry
	Deregister(serviceName string) error
	// GetService Resolves the servicename and returns the service details
	GetService(serviceName string) ([]Entry, error)
	// ListServices returns all the services in the registry
	ListServices() ([]Entry, error)
	// Watch sets the registry to watch mode so that it tracks any updates
	Watch() error
	// CancelWatch stops the registry watch mode
	CancelWatch() error
	// String returns the description of the registry
	String() string
}
