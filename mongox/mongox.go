package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database is the mongox database interface
type Database interface {
	Client() *mongo.Client
	Context() context.Context
	Name() string
	New(ctx context.Context) Database
	GetCollectionOf(document interface{}) *mongo.Collection
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
