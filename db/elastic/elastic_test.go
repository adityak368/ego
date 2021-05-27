package elastic

import (
	"testing"

	"github.com/adityak368/ego/db"
	"github.com/adityak368/swissknife/logger/v2"
	"github.com/stretchr/testify/require"
)

// Test Model
type Test struct {
	Name string `json:"name" validate:"required"`
}

// CreateIndexes creates the indexes for a model
func (m *Test) CreateIndexes(db db.Database) error {

	r, err := db.(*DB).InsertDocument("test", m)
	if err != nil {
		return err
	}

	logger.Print(r)

	return nil
}

// String returns the string representation of the model
func (m *Test) String() string {
	return "Test"
}

func TestElastic(t *testing.T) {

	r := require.New(t)

	ElasticDB := New()
	ElasticDB.Init(db.Options{
		Name:    "ElasdticDB",
		Address: "http://localhost:9200",
	})
	err := ElasticDB.Connect()
	r.Nil(err)

	m := &Test{}

	err = m.CreateIndexes(ElasticDB)
	r.Nil(err)

	r.NotPanics(func() { ElasticDB.(*DB).PrintIndexes() })
}
