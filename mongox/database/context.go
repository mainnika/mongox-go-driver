package database

import (
	"context"
)

type ctxDatabaseKey struct{}

// GetFromContext function extracts the request data from context
func GetFromContext(ctx context.Context) (q *Database, ok bool) {
	q, ok = ctx.Value(ctxDatabaseKey{}).(*Database)
	if !ok {
		return nil, false
	}

	return q, true
}

// WithContext creates the new context with a database attached
func WithContext(ctx context.Context, q *Database) (withQuery context.Context) {
	db := NewDatabase(ctx, q.Client(), q.Name())
	return context.WithValue(ctx, ctxDatabaseKey{}, db)
}
