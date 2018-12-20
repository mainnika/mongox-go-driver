package mongox

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Saver interface {
	Save(db *Database) error
}

type Deleter interface {
	Delete(db *Database) error
}

type Loader interface {
	Load(db *Database, filters ...interface{}) error
}

type BaseObjectID interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
}
