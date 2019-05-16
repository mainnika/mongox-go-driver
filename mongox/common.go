package mongox

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Saver is an interface for documents that can be saved
type Saver interface {
	Save(db *Database) error
}

// Deleter is an interface for documents that can be deleted
type Deleter interface {
	Delete(db *Database) error
}

// Loader is an interface for documents that can be loaded
type Loader interface {
	Load(db *Database, filters ...interface{}) error
}

// Resetter is an interface for documenta that can be resetted
type Resetter interface {
	Reset()
}

// BaseObjectID is an interface for documents that have objectId type for the _id field
type BaseObjectID interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
}

// BaseString is an interface for documents that have string type for the _id field
type BaseString interface {
	GetID() string
	SetID(id string)
}

// BaseObject is an interface for documents that have object type for the _id field
type BaseObject interface {
	GetID() primitive.D
	SetID(id primitive.D)
}
