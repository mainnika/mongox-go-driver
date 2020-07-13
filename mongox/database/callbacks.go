package database

import (
	"context"

	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

func onDecode(ctx context.Context, iter interface{}, callbacks ...query.OnDecode) (err error) {

	for _, cb := range callbacks {
		err = cb(ctx, iter)
		if err != nil {
			return
		}
	}

	return
}
