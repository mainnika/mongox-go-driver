package database

import (
	"context"
	"math/rand"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/database"
)

// EphemeralDatabase is a temporary database connection that will be destroyed after close
type EphemeralDatabase struct {
	mongox.Database
}

// NewEphemeral creates new mongo connection
func NewEphemeral(URI string) (db *EphemeralDatabase, err error) {

	name := strconv.Itoa(rand.Int())
	opts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.Background(), opts)

	db = &EphemeralDatabase{Database: database.NewDatabase(client, name)}

	return
}

// Close the connection and drop database
func (e *EphemeralDatabase) Close() error {
	return e.Client().Database(e.Name()).Drop(e.Context())
}
