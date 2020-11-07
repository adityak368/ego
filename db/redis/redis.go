package redis

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/adityak368/ego/db"
	"github.com/adityak368/swissknife/logger"
	"github.com/go-redis/redis"
)

// DB implements redis database
type DB struct {
	options db.Options
	client  *redis.Client
}

// Address is the bind address
func (r *DB) Address() string {
	return r.options.Address
}

// Init initializes the db connection
func (r *DB) Init(opts db.Options) error {
	r.options = opts
	return nil
}

// Options Returns the db options
func (r *DB) Options() db.Options {
	return r.options
}

// Connect connects to the db
func (r *DB) Connect(handlers ...db.Handler) error {

	database, err := strconv.Atoi(r.options.Database)
	if err != nil {
		return err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     r.Address(),
		Password: r.options.Password,
		DB:       database,
	})

	if _, err := client.Ping().Result(); err != nil {
		return err
	}

	r.client = client

	for _, fn := range handlers {
		err := fn(r)
		if err != nil {
			return err
		}
	}

	logger.Info(r)

	return nil
}

// Disconnect disconnects the db connection
func (r *DB) Disconnect(handlers ...db.Handler) error {

	if r.client == nil {
		return errors.New("[DB]: Not connected to Redis")
	}

	err := r.client.Close()
	if err != nil {
		return err
	}

	for _, fn := range handlers {
		err := fn(r)
		if err != nil {
			return err
		}
	}
	return err
}

// String returns the description of the database
func (r *DB) String() string {
	return fmt.Sprintf("[DB]: Redis connected to %s", r.Address())
}

// Handle returns the raw connection handle to the db
func (r *DB) Handle() interface{} {
	return r.client
}

// New returns a new grpc server
func New() db.Database {
	return &DB{}
}
