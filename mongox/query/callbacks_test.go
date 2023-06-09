package query_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

func TestCallbacks_InvokeOk(t *testing.T) {

	mocked := mock.Mock{}

	var callbacks = query.Callbacks{
		func(ctx context.Context, iter interface{}) (err error) {
			return mocked.Called(ctx, iter).Error(0)
		},
		func(ctx context.Context, iter interface{}) (err error) {
			return mocked.Called(ctx, iter).Error(0)
		},
	}

	ctx := context.Background()
	iter := int64(42)

	mocked.On("func1", ctx, iter).Return(nil).Once()
	mocked.On("func2", ctx, iter).Return(nil).Once()

	assert.NoError(t, callbacks.Invoke(ctx, iter))
	assert.True(t, mocked.AssertExpectations(t))
}

func TestCallbacks_InvokeStopIfError(t *testing.T) {

	mocked := mock.Mock{}

	var callbacks = query.Callbacks{
		func(ctx context.Context, iter interface{}) (err error) {
			return mocked.Called(ctx, iter).Error(0)
		},
		func(ctx context.Context, iter interface{}) (err error) {
			t.FailNow()
			return
		},
	}

	ctx := context.Background()
	iter := int(42)

	mocked.On("func1", ctx, iter).Return(fmt.Errorf("wat"))

	assert.EqualError(t, callbacks.Invoke(ctx, iter), "wat")
	assert.True(t, mocked.AssertExpectations(t))
}
