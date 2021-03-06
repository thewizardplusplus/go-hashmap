package hashmap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockKeyWithID struct {
	MockKey

	ID int
}

func NewMockKeyWithID(id int) *MockKeyWithID {
	return &MockKeyWithID{ID: id}
}

func TestWithInterruption(test *testing.T) {
	type wrapperArgs struct {
		ctx     context.Context
		handler HandlerInterface
	}
	type handlerArgs struct {
		key   Key
		value interface{}
	}

	for _, data := range []struct {
		name        string
		wrapperArgs wrapperArgs
		handlerArgs handlerArgs
		want        assert.BoolAssertionFunc
	}{
		{
			name: "not interrupted and positive",
			wrapperArgs: wrapperArgs{
				ctx: context.Background(),
				handler: func() HandlerInterface {
					handler := new(MockHandlerInterface)
					handler.On("Handle", NewMockKeyWithID(23), "data").Return(true)

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				key:   NewMockKeyWithID(23),
				value: "data",
			},
			want: assert.True,
		},
		{
			name: "not interrupted and negative",
			wrapperArgs: wrapperArgs{
				ctx: context.Background(),
				handler: func() HandlerInterface {
					handler := new(MockHandlerInterface)
					handler.On("Handle", NewMockKeyWithID(23), "data").Return(false)

					return handler
				}(),
			},
			handlerArgs: handlerArgs{
				key:   NewMockKeyWithID(23),
				value: "data",
			},
			want: assert.False,
		},
		{
			name: "interrupted",
			wrapperArgs: wrapperArgs{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				}(),
				handler: new(MockHandlerInterface),
			},
			handlerArgs: handlerArgs{
				key:   NewMockKeyWithID(23),
				value: "data",
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := WithInterruption(
				data.wrapperArgs.ctx,
				data.wrapperArgs.handler.Handle,
			)
			require.NotNil(test, handler)

			got := handler(data.handlerArgs.key, data.handlerArgs.value)

			mock.AssertExpectationsForObjects(
				test,
				data.wrapperArgs.handler,
				data.handlerArgs.key,
			)
			data.want(test, got)
		})
	}
}
