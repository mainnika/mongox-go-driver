package common

import (
	"context"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// ManyLoader is a controller for a database cursor
type ManyLoader struct {
	mongo.Cursor
	ctx    context.Context
	target interface{}
}

// Get loads documents to a target or returns an error
func (l *ManyLoader) Get() error {

	hasNext := l.Next(l.ctx)

	if !hasNext {
		return errors.NotFoundErrorf("%s", mongo.ErrNoDocuments)
	}

	err := l.Decode(l.target)
	if err != nil {
		return errors.InternalErrorf("can't decode desult: %s", err)
	}

	return nil
}

// Close cursor
func (l *ManyLoader) Close() error {

	return l.Cursor.Close(l.ctx)
}

// LoadMany function loads documents one by one into a target channel
func LoadMany(db *mongox.Database, target interface{}, composed *query.Query) (*ManyLoader, error) {

	collection := db.GetCollectionOf(target)
	opts := options.Find()

	opts.Sort = composed.Sorter()
	opts.Limit = composed.Limiter()

	cursor, err := collection.Find(db.Context(), composed.M(), opts)
	if err != nil {
		return nil, errors.InternalErrorf("can't create find result: %s", err)
	}

	l := &ManyLoader{Cursor: cursor, ctx: db.Context(), target: target}

	return l, nil
}
