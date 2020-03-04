package tempdb

import (
	"context"
	"math/rand"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/database"
)

// TempDB is a temporary database connection that will be destroyed after close
type TempDB struct {
	mongox.Database
}

// NewTempDB creates new mongo connection
func NewTempDB(URI string) (tempdb *TempDB, err error) {

	name := strconv.Itoa(rand.Int())
	opts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.Background(), opts)

	tempdb = &TempDB{Database: database.NewDatabase(client, name)}

	return
}

// Close the connection and drop database
func (tdb *TempDB) Close() {
	_ = tdb.Client().Database(tdb.Name()).Drop(tdb.Context())
}
