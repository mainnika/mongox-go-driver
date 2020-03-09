package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database is the mongox database interface
type Database interface {
	Client() MongoClient
	Context() context.Context
	Name() string
	New(ctx context.Context) Database
	GetCollectionOf(document interface{}) MongoCollection
	Count(target interface{}, filters ...interface{}) (int64, error)
	DeleteArray(target interface{}) error
	DeleteOne(target interface{}, filters ...interface{}) error
	LoadArray(target interface{}, filters ...interface{}) error
	LoadOne(target interface{}, filters ...interface{}) error
	LoadStream(target interface{}, filters ...interface{}) (StreamLoader, error)
	SaveOne(source interface{}) error
}

// StreamLoader is a interface to control database cursor
type StreamLoader interface {
	DecodeNext() error
	Decode() error
	Next() error
	Close() error
	Err() error
}

// MongoClient is the mongo client interface
type MongoClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error)
	ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	NumberSessionsInProgress() int
}

// MongoCollection is the mongo collection interface
type MongoCollection interface {
	Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error)
	Name() string
	Database() *mongo.Database
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	Indexes() mongo.IndexView
	Drop(ctx context.Context) error
}

// Saver is an interface for documents that can be saved
type Saver interface {
	Save(db Database) error
}

// Deleter is an interface for documents that can be deleted
type Deleter interface {
	Delete(db Database) error
}

// Loader is an interface for documents that can be loaded
type Loader interface {
	Load(db Database, filters ...interface{}) error
}

// Resetter is an interface for documenta that can be resetted
type Resetter interface {
	Reset()
}

// OIDBased is an interface for documents that have objectId type for the _id field
type OIDBased interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
}

// StringBased is an interface for documents that have string type for the _id field
type StringBased interface {
	GetID() string
	SetID(id string)
}

// JSONBased is an interface for documents that have object type for the _id field
type JSONBased interface {
	GetID() primitive.D
	SetID(id primitive.D)
}

type InterfaceBased interface {
	GetID() interface{}
	SetID(id interface{})
}
