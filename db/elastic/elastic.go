package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/adityak368/ego/db"
	"github.com/adityak368/swissknife/logger/v2"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// DB implements elastic database
type DB struct {
	options db.Options
	client  *elasticsearch.Client
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

	cfg := elasticsearch.Config{
		Addresses: []string{
			m.options.Address,
		},
	}

	if m.options.Username != "" {
		cfg.Username = m.options.Username
		cfg.Password = m.options.Password
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return err
	}

	res, err := es.Info()
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.String())
	}

	defer res.Body.Close()

	m.client = es

	for _, fn := range handlers {
		err := fn(m)
		if err != nil {
			return err
		}
	}

	var r map[string]interface{}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	// Print client and server version numbers.
	logger.Info().Msgf("[DB]: Client: %s", elasticsearch.Version)
	logger.Info().Msgf("[DB]: Server: %s", r["version"].(map[string]interface{})["number"])

	logger.Info().Msg(m.String())

	return nil
}

// Disconnect disconnects the db connection
func (m *DB) Disconnect(handlers ...db.Handler) error {

	if m.client == nil {
		return errors.New("[DB]: Cannot Disconnect. Not connected to elastic")
	}

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
	return fmt.Sprintf("[DB]: Elastic Connected to %s", m.Address())
}

// Handle returns the raw connection handle to the db
func (m *DB) Handle() interface{} {
	return m.client
}

// PrintIndexes prints all the indexes in the db
func (m *DB) PrintIndexes() {
	if m.client == nil {
		logger.Warn().Msg("[DB]: Cannot print indexes. Not connected to elastic")
		return
	}

	req := esapi.CatIndicesRequest{
		Pretty: true,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		logger.Warn().Err(err).Msg("[DB]: Cannot print indexes. Not connected to elastic")
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error().Msg("[DB]: Something went wrong listing indexes")
		return
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		logger.Error().Msg("[DB]: Could not read index response")
		return
	}

	logger.Info().Msgf("[DB]: %s", buf)
}

// InsertDocument inserts a document to an index
func (m *DB) InsertDocument(index string, document interface{}) (map[string]interface{}, error) {

	b, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(string(b)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

// InsertDocumentAsString inserts a document to an index
func (m *DB) InsertDocumentAsString(index string, data string) (map[string]interface{}, error) {

	req := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(string(data)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

// New returns a new grpc server
func New() db.Database {
	return &DB{}
}
