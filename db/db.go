package db

// Database defines the interface for connection with database
type Database interface {
	// Init initializes the db connection
	Init(opts Options) error
	// Options Returns the db options
	Options() Options
	// Address Returns the db bind interface
	Address() string
	// Connect connects to the db
	Connect(handlers ...Handler) error
	// Disconnect disconnects from the db
	Disconnect(handlers ...Handler) error
	// Handle returns the raw connection handle to the db
	Handle() interface{}
	// String returns the description of the database
	String() string
}

// Handler is used to run a function on db events
type Handler func(Database) error
