package query

import (
	"context"
)

type ctxQueryKey struct{}

// GetFromContext function extracts the request data from context
func GetFromContext(ctx context.Context) (q *Query, ok bool) {
	q, ok = ctx.Value(ctxQueryKey{}).(*Query)
	if !ok {
		return nil, false
	}

	return q, true
}

// WithContext function creates the new context with request data
func WithContext(ctx context.Context, q *Query) (withQuery context.Context) {
	return context.WithValue(ctx, ctxQueryKey{}, q)
}
