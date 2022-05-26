package elastic

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/adityak368/ego/db"
	"github.com/adityak368/swissknife/logger/v2"
	"github.com/gocql/gocql"
)

// DB implements cassandra database
type DB struct {
	options db.Options
	client  *gocql.Session
}

// Address is the bind address
func (m *DB) Address() string {
	return m.options.Address
}

// Init initializes the db connection
func (m *DB) Init(opts db.Options) error {
	m.options = opts
	return nil
}

// Options Returns the db options
func (m *DB) Options() db.Options {
	return m.options
}

// Connect connects to the db
func (m *DB) Connect(handlers ...db.Handler) error {

	cluster := gocql.NewCluster(m.options.Address)
	cluster.Keyspace = m.options.Database
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 13

	if m.options.Username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: m.options.Username,
			Password: m.options.Password,
		}
	}

	session, err := cluster.CreateSession()

	if err != nil {
		return err
	}

	m.client = session

	for _, fn := range handlers {
		err := fn(m)
		if err != nil {
			return err
		}
	}

	logger.Info().Msg(m.String())

	return nil

}

// Disconnect disconnects the db connection
func (m *DB) Disconnect(handlers ...db.Handler) error {

	if m.client == nil || m.client.Closed() {
		return errors.New("[DB]: Cannot Disconnect. Not connected to cassandra")
	}

	m.client.Close()

	for _, fn := range handlers {
		err := fn(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// String returns the description of the database
func (m *DB) String() string {
	return fmt.Sprintf("[DB]: Cassandra Connected to %s", m.Address())
}

// Handle returns the raw connection handle to the db
func (m *DB) Handle() interface{} {
	return m.client
}

// New returns a new grpc server
func New() db.Database {
	return &DB{}
}
