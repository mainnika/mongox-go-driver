package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/database"
)

// defaultURI is a mongodb uri that is being used by tests
var defaultURI = "mongodb://localhost"

// EphemeralDatabase is a temporary database connection that will be destroyed after close
type EphemeralDatabase struct {
	mongox.Database
}

func init() {
	envURI := os.Getenv("MONGODB_URI")
	if envURI != "" {
		defaultURI = envURI
	}
}

// NewEphemeral creates new mongo connection
func NewEphemeral(URI string) (db *EphemeralDatabase, err error) {
	return NewEphemeralWithContext(context.Background(), URI)
}

func NewEphemeralWithContext(ctx context.Context, URI string) (db *EphemeralDatabase, err error) {
	if URI == "" {
		URI = defaultURI
	}

	name := primitive.NewObjectID().Hex()
	opts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	db = &EphemeralDatabase{Database: database.NewDatabase(ctx, client, name)}

	return db, nil
}

// Close the connection and drop database
func (e *EphemeralDatabase) Close() (err error) {
	return e.Client().Database(e.Name()).Drop(e.Context())
}
