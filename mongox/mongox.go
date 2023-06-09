package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Reexport basic mongo structs
type (
	Cursor     = mongo.Cursor
	Client     = mongo.Client
	Collection = mongo.Collection
)

// Database is the mongox database interface
type Database interface {
	Client() (client *Client)
	Context() (context context.Context)
	Name() (name string)
	GetCollectionOf(document interface{}) (collection *Collection, err error)
	Count(target interface{}, filters ...interface{}) (count int64, err error)
	DeleteArray(target interface{}, filters ...interface{}) (err error)
	DeleteOne(target interface{}, filters ...interface{}) (err error)
	LoadArray(target interface{}, filters ...interface{}) (err error)
	LoadOne(target interface{}, filters ...interface{}) (err error)
	LoadStream(target interface{}, filters ...interface{}) (loader StreamLoader, err error)
	SaveOne(source interface{}, filters ...interface{}) (err error)
	UpdateOne(target interface{}, filters ...interface{}) (err error)
	IndexEnsure(cfg interface{}, document interface{}) (err error)
}

// StreamLoader is a interface to control database cursor
type StreamLoader interface {
	Cursor() (cursor *Cursor)
	DecodeNextMsg(i interface{}) (err error)
	DecodeMsg(i interface{}) (err error)
	Next() (err error)
	Close() (err error)
	Err() (err error)
}

// OIDBased is an interface for documents that have objectId type for the _id field
type OIDBased interface {
	GetID() (id primitive.ObjectID)
	SetID(id primitive.ObjectID)
}

// StringBased is an interface for documents that have string type for the _id field
type StringBased interface {
	GetID() (id string)
	SetID(id string)
}

// DocBased is an interface for documents that have object type for the _id field
type DocBased interface {
	GetID() (id primitive.D)
	SetID(id primitive.D)
}

// InterfaceBased is an interface for documents that have custom declated type for the _id field
type InterfaceBased interface {
	GetID() (id interface{})
	SetID(id interface{})
}
