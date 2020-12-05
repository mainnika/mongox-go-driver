package database

import (
	"context"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// StreamLoader is a controller for a database cursor
type StreamLoader struct {
	cur    *mongox.Cursor
	query  *query.Query
	ctx    context.Context
	target interface{}
}

// DecodeNext loads next documents to a target or returns an error
func (l *StreamLoader) DecodeNext() (err error) {

	err = l.Next()
	if err != nil {
		return
	}

	err = l.Decode()
	if err != nil {
		return
	}

	return
}

// Decode function decodes the current cursor document into the target
func (l *StreamLoader) Decode() (err error) {

	if created := base.Reset(l.target); created {
		err = l.query.OnDecode().Invoke(l.ctx, l.target)
	}
	if err != nil {
		return
	}

	err = l.cur.Decode(l.target)
	if err != nil {
		return
	}

	err = l.query.OnDecode().Invoke(l.ctx, l.target)
	if err != nil {
		return
	}

	return
}

// Next loads next documents but doesn't perform decoding
func (l *StreamLoader) Next() (err error) {

	hasNext := l.cur.Next(l.ctx)
	err = l.cur.Err()

	if err != nil {
		return
	}
	if !hasNext {
		err = mongox.ErrNoDocuments
	}

	return
}

func (l *StreamLoader) Cursor() (cursor *mongox.Cursor) {
	return l.cur
}

// Close cursor
func (l *StreamLoader) Close() (err error) {

	closerr := l.cur.Close(l.ctx)
	invokerr := l.query.OnClose().Invoke(l.ctx, l.target)

	if closerr != nil {
		err = closerr
		return
	}

	if invokerr != nil {
		err = invokerr
		return
	}

	return
}

func (l *StreamLoader) Err() (err error) {
	return l.cur.Err()
}
