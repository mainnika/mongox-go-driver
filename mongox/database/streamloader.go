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
		_ = l.query.OnDecode().Invoke(l.ctx, i)
	}
	if err != nil {
		return err
	}

	err = l.cur.Decode(i)
	if err != nil {
		return err
	}

	_ = l.query.OnDecode().Invoke(l.ctx, i)

	return nil
}

// Next loads next documents but doesn't perform decoding
func (l *StreamLoader) Next() (err error) {
	hasNext := l.cur.Next(l.ctx)
	err = l.cur.Err()
	if err != nil {
		return err
	}
	if !hasNext {
		return mongox.ErrNoDocuments
	}

	return nil
}

// Cursor returns the underlying cursor
func (l *StreamLoader) Cursor() (cursor *mongox.Cursor) {
	return l.cur
}

// Close stream loader and the underlying cursor
func (l *StreamLoader) Close() (err error) {
	defer func() { _ = l.query.OnClose().Invoke(l.ctx, nil) }()

	err = l.cur.Close(l.ctx)
	if err != nil {
		return err
	}

	return nil
}

// Err returns the last error
func (l *StreamLoader) Err() (err error) {
	return l.cur.Err()
}

func (l *StreamLoader) Context() (ctx context.Context) {
	return l.ctx
}
