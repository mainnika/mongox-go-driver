package database

import (
	"context"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/base"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// StreamLoader is a controller for a database cursor
type StreamLoader struct {
	cur   *mongox.Cursor
	query *query.Query
	ctx   context.Context
	ref   interface{}
}

// DecodeNextMsg decodes the next document to an interface or returns an error
func (l *StreamLoader) DecodeNextMsg(i interface{}) (err error) {

	err = l.Next()
	if err != nil {
		return
	}

	err = l.DecodeMsg(i)
	if err != nil {
		return
	}

	return
}

// DecodeMsg decodes the current cursor document into an interface
func (l *StreamLoader) DecodeMsg(i interface{}) (err error) {

	if created := base.Reset(i); created {
		err = l.query.OnDecode().Invoke(l.ctx, i)
	}
	if err != nil {
		return
	}

	err = l.cur.Decode(i)
	if err != nil {
		return
	}

	err = l.query.OnDecode().Invoke(l.ctx, i)
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
	invokerr := l.query.OnClose().Invoke(l.ctx, l.ref)

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
