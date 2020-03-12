package database

import (
	"context"
	"fmt"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
)

// StreamLoader is a controller for a database cursor
type StreamLoader struct {
	cur    *mongox.Cursor
	ctx    context.Context
	target interface{}
}

// DecodeNext loads next documents to a target or returns an error
func (l *StreamLoader) DecodeNext() error {

	hasNext := l.cur.Next(l.ctx)

	if l.cur.Err() != nil {
		return l.cur.Err()
	}
	if !hasNext {
		return mongox.ErrNoDocuments
	}

	base.Reset(l.target)

	err := l.cur.Decode(l.target)
	if err != nil {
		return fmt.Errorf("can't decode desult: %w", err)
	}

	return nil
}

// Decode function decodes the current cursor document into the target
func (l *StreamLoader) Decode() error {

	base.Reset(l.target)

	err := l.cur.Decode(l.target)
	if err != nil {
		return fmt.Errorf("can't decode desult: %w", err)
	}

	return nil
}

// Next loads next documents but doesn't perform decoding
func (l *StreamLoader) Next() error {

	hasNext := l.cur.Next(l.ctx)

	if l.cur.Err() != nil {
		return l.cur.Err()
	}
	if !hasNext {
		return mongox.ErrNoDocuments
	}

	return nil
}

func (l *StreamLoader) Cursor() *mongox.Cursor {
	return l.cur
}

// Close cursor
func (l *StreamLoader) Close() error {
	return l.cur.Close(l.ctx)
}

func (l *StreamLoader) Err() error {
	return l.cur.Err()
}
