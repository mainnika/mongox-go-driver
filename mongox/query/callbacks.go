package query

import (
	"context"
)

type Callback func(ctx context.Context, iter interface{}) (err error)
type Callbacks []Callback

type (
	OnDecode Callback
	OnClose  Callback
)

// Invoke callbacks sequence
func (c Callbacks) Invoke(ctx context.Context, iter interface{}) (err error) {

	for _, cb := range c {
		err = cb(ctx, iter)
		if err != nil {
			return
		}
	}

	return
}
