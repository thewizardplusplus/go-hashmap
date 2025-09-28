package hashmap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithInterruption(test *testing.T) {
	type wrapperArgs struct {
		ctx     context.Context
		handler func(test *testing.T) handlerInterface
	}
	type handlerArgs struct {
		key   func(test *testing.T) Key
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
				handler: func(test *testing.T) handlerInterface {
					handler := NewMockhandlerInterface(test)
					handler.EXPECT().Handle(NewMockKey(test), "data").Return(true)

					return handler
				},
			},
			handlerArgs: handlerArgs{
				key:   func(test *testing.T) Key { return NewMockKey(test) },
				value: "data",
			},
			want: assert.True,
		},
		{
			name: "not interrupted and negative",
			wrapperArgs: wrapperArgs{
				ctx: context.Background(),
				handler: func(test *testing.T) handlerInterface {
					handler := NewMockhandlerInterface(test)
					handler.EXPECT().Handle(NewMockKey(test), "data").Return(false)

					return handler
				},
			},
			handlerArgs: handlerArgs{
				key:   func(test *testing.T) Key { return NewMockKey(test) },
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
				handler: func(test *testing.T) handlerInterface {
					return NewMockhandlerInterface(test)
				},
			},
			handlerArgs: handlerArgs{
				key:   func(test *testing.T) Key { return NewMockKey(test) },
				value: "data",
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := WithInterruption(
				data.wrapperArgs.ctx,
				data.wrapperArgs.handler(test).Handle,
			)
			require.NotNil(test, handler)

			got := handler(data.handlerArgs.key(test), data.handlerArgs.value)

			data.want(test, got)
		})
	}
}
