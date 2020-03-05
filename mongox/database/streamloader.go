package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
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

	err := l.Cursor.Decode(l.target)
	if err != nil {
		return fmt.Errorf("can't decode desult: %w", err)
	}

	return nil
}

func (l *StreamLoader) Decode() error {

	base.Reset(l.target)

	err := l.Cursor.Decode(l.target)
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
