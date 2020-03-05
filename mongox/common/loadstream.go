package common

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
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

	if l.Cursor.Err() != nil {
		return l.Cursor.Err()
	}
	if !hasNext {
		return mongo.ErrNoDocuments
	}

	base.Reset(l.target)

	err := l.Decode(l.target)
	if err != nil {
		return fmt.Errorf("can't decode desult: %w", err)
	}

	return nil
}

// Next loads next documents but doesn't perform decoding
func (l *StreamLoader) Next() error {

	hasNext := l.Cursor.Next(l.ctx)

	if l.Cursor.Err() != nil {
		return l.Cursor.Err()
	}
	if !hasNext {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Close cursor
func (l *StreamLoader) Close() error {

	return l.Cursor.Close(l.ctx)
}

// LoadStream function loads documents one by one into a target channel
func LoadStream(db mongox.Database, target interface{}, filters ...interface{}) (*StreamLoader, error) {

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
		return nil, fmt.Errorf("can't create find result: %w", err)
	}

	l := &StreamLoader{Cursor: cursor, ctx: db.Context(), target: target}

	return l, nil
}
