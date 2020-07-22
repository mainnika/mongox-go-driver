package query

import (
	"context"
)

type ctxQueryKey struct{}

// GetFromContext function extracts the request data from context
func GetFromContext(ctx context.Context) (q *Query, ok bool) {
	q, ok = ctx.Value(ctxQueryKey{}).(*Query)
	return
}

// WithContext function creates the new context with request data
func WithContext(ctx context.Context, q *Query) (withQuery context.Context) {
	withQuery = context.WithValue(ctx, ctxQueryKey{}, q)
	return
}
