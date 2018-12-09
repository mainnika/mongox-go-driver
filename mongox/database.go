package mongox

import (
	"context"
	"reflect"

	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Database handler
type Database struct {
	client *mongo.Client
	dbname string
	ctx    context.Context
}

// NewDatabase function creates new database instance with mongo client and empty context
func NewDatabase(client *mongo.Client, dbname string) *Database {

	db := &Database{}
	db.client = client
	db.dbname = dbname

	return db
}

// Client function returns a mongo client
func (d *Database) Client() *mongo.Client {
	return d.client
}

// Context function returns a context
func (d *Database) Context() context.Context {
	return d.ctx
}

// Name function returns a database name
func (d *Database) Name() string {
	return d.dbname
}

// New function creates new database context with same client
func (d *Database) New(ctx context.Context) *Database {

	if ctx != nil {
		ctx = context.Background()
	}

	return &Database{
		client: d.client,
		dbname: d.dbname,
		ctx:    ctx,
	}
}

// GetCollectionOf returns the collection object by the «collection» tag of the given document;
// the «collection» tag should exists, e.g.:
// type Foobar struct {
//     base.ObjectID `bson:",inline" json:",inline" collection:"foobars"`
// 	   ...
// Will panic if there is no «collection» tag
func (d *Database) GetCollectionOf(document interface{}) *mongo.Collection {

	el := reflect.TypeOf(document).Elem()
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		tag := field.Tag
		found, ok := tag.Lookup("collection")
		if !ok {
			continue
		}

		return d.client.Database(d.dbname).Collection(found)
	}

	panic(errors.InternalErrorf("document %v does not have a collection tag", document))
}
