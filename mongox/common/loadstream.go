package common

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/base"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
)

// StreamLoader is a controller for a database cursor
type StreamLoader struct {
	*mongo.Cursor
	ctx    context.Context
	target interface{}
}

// DecodeNext loads next documents to a target or returns an error
func (l *StreamLoader) DecodeNext() error {

	hasNext := l.Cursor.Next(l.ctx)

	if !hasNext {
		return errors.NotFoundErrorf("%s", mongo.ErrNoDocuments)
	}

	base.Reset(l.target)

	err := l.Decode(l.target)
	if err != nil {
		return errors.InternalErrorf("can't decode desult: %s", err)
	}

	return nil
}

// Next loads next documents but doesn't perform decoding
func (l *StreamLoader) Next() error {

	hasNext := l.Cursor.Next(l.ctx)
	if !hasNext {
		return errors.NotFoundErrorf("%s", mongo.ErrNoDocuments)
	}

	return nil
}

// Close cursor
func (l *StreamLoader) Close() error {

	return l.Cursor.Close(l.ctx)
}

// LoadStream function loads documents one by one into a target channel
func LoadStream(db *mongox.Database, target interface{}, filters ...interface{}) (*StreamLoader, error) {

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

	l := &StreamLoader{Cursor: cursor, ctx: db.Context(), target: target}

	return l, nil
}
