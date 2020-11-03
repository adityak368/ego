package db

// Model defines the interface for a model that can saved in a db
type Model interface {
	// CreateIndexes creates the required indexes
	CreateIndexes(db Database) error
	// PrintIndexes prints the indexes of the model
	PrintIndexes(db Database)
	// String returns the description of the model
	String() string
}
