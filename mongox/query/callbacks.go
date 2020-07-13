package query

import (
	"context"
)

type OnDecode func(ctx context.Context, iter interface{}) (err error)
