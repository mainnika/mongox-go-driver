package database

import (
	"context"
	"reflect"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
)

// Database handler
type Database struct {
	client *mongox.Client
	name   string
	ctx    context.Context
}

// NewDatabase function creates new database instance with mongo client and empty context
func NewDatabase(ctx context.Context, client *mongox.Client, name string) (db mongox.Database) {
	db = &Database{
		client: client,
		name:   name,
		ctx:    ctx,
	}

	return db
}

// Client function returns a mongo client
func (d *Database) Client() (client *mongox.Client) {
	return d.client
}

// Name function returns a database name
func (d *Database) Name() (name string) {
	return d.name
}

// Context function returns a context
func (d *Database) Context() (ctx context.Context) {
	ctx = d.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

// GetCollectionOf returns the collection object by the «collection» tag of the given document;
//
//	 example:
//		type Foobar struct {
//		    base.ObjectID `bson:",inline" json:",inline" collection:"foobars"`
//			   ...
func (d *Database) GetCollectionOf(document interface{}) (collection *mongox.Collection, err error) {
	el := reflect.TypeOf(document).Elem()
	numField := el.NumField()
	databaseName := d.name

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		tag := field.Tag
		collectionName, found := tag.Lookup("collection")
		if !found {
			continue
		}

		return d.client.Database(databaseName).Collection(collectionName), nil
	}

	return nil, mongox.ErrNoCollection
}
