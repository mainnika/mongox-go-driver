package common

import (
	"context"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"go.mongodb.org/mongo-driver/mongo"
)

// ManyLoader is a controller for a database cursor
type ManyLoader struct {
	*mongo.Cursor
	ctx    context.Context
	target interface{}
}

// GetNext loads next documents to a target or returns an error
func (l *ManyLoader) GetNext() error {

	hasNext := l.Cursor.Next(l.ctx)

	if !hasNext {
		return errors.NotFoundErrorf("%s", mongo.ErrNoDocuments)
	}

	resettable, canReset := l.target.(mongox.Resetter)
	if canReset {
		resettable.Reset()
	}

	err := l.Decode(l.target)
	if err != nil {
		return errors.InternalErrorf("can't decode desult: %s", err)
	}

	return nil
}

// Next loads next documents but doesn't perform decoding
func (l *ManyLoader) Next() error {

	hasNext := l.Cursor.Next(l.ctx)
	if !hasNext {
		return errors.NotFoundErrorf("%s", mongo.ErrNoDocuments)
	}

	return nil
}

// Close cursor
func (l *ManyLoader) Close() error {

	return l.Cursor.Close(l.ctx)
}

// LoadMany function loads documents one by one into a target channel
func LoadMany(db *mongox.Database, target interface{}, filters ...interface{}) (*ManyLoader, error) {

	var cursor *mongo.Cursor
	var err error

	composed := query.Compose(filters...)
	hasPreloader, _ := composed.Preloader()

	if hasPreloader {
		cursor, err = createAggregateLoad(db, target, composed)
	} else {
		cursor, err = createSimpleLoad(db, target, composed)
	}
	if err != nil {
		return nil, errors.InternalErrorf("can't create find result: %s", err)
	}

	l := &ManyLoader{Cursor: cursor, ctx: db.Context(), target: target}

	return l, nil
}
