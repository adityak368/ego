package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adityak368/ego/db"
	"github.com/adityak368/swissknife/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB implements mongodb database
type DB struct {
	options    db.Options
	client     *mongo.Client
	connection *mongo.Database
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

	//certPath, isCertificateAvailable := os.LookupEnv("MONGO_DB_CERT_PATH")

	clientOptions := &options.ClientOptions{
		AppName: &m.options.Name,
		Hosts:   []string{m.Address()},
	}
	if m.options.Username != "" {
		clientOptions.Auth = &options.Credential{
			Username: m.options.Username,
			Password: m.options.Password,
		}
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	m.client = client
	m.connection = client.Database(m.options.Database)

	for _, fn := range handlers {
		err := fn(m)
		if err != nil {
			return err
		}
	}

	logger.Info(m)

	return nil
}

// Disconnect disconnects the db connection
func (m *DB) Disconnect(handlers ...db.Handler) error {

	if m.client == nil {
		return errors.New("[DB]: Cannot Disconnect. Not connected to MongoDB")
	}

	err := m.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	for _, fn := range handlers {
		err := fn(m)
		if err != nil {
			return err
		}
	}

	return err
}

// String returns the description of the database
func (m *DB) String() string {
	return fmt.Sprintf("[DB]: MongoDB Connected to %s", m.Address())
}

// Handle returns the raw connection handle to the db
func (m *DB) Handle() interface{} {
	return m.connection
}

// PrintIndexes prints all the indexes in the db
func (m *DB) PrintIndexes(collection string) {
	if m.client == nil {
		logger.Warn("[DB]: Cannot print indexes. Not connected to MongoDB")
		return
	}

	if m.connection == nil {
		logger.Warn("[DB]: Cannot print indexes. Not a valid database")
		return
	}

	c := m.connection.Collection(collection)
	duration := 10 * time.Second
	batchSize := int32(10)
	cur, err := c.Indexes().List(context.Background(), &options.ListIndexesOptions{BatchSize: &batchSize, MaxTime: &duration})
	if err != nil {
		logger.Error("[DB]: Something went wrong listing indexes", err)
	}
	for cur.Next(context.Background()) {
		index := bson.D{}
		cur.Decode(&index)
		logger.Infof("Index found: %v", index)
	}
}

// New returns a new grpc server
func New() db.Database {
	return &DB{}
}
