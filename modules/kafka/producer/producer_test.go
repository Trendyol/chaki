package producer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

type testProducer struct {
	msg       string
	returnErr bool
}

func (p *testProducer) Produce(ctx context.Context, msg ...Message) error {
	if p.returnErr {
		return errTest
	}
	p.msg = string(msg[0].Value)
	return nil
}

func Test_buildInterceptors(t *testing.T) {
	t.Run("it should build interceptors when interceptors provided", func(t *testing.T) {
		// Given
		var (
			p                = &testProducer{}
			interceptorCalls []int
			interceptors     = []Interceptor{
				InterceptorFunc(func(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
					interceptorCalls = append(interceptorCalls, 1)
					return next(ctx, msgs...)
				}),
				InterceptorFunc(func(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
					interceptorCalls = append(interceptorCalls, 2)
					return next(ctx, msgs...)
				}),
				InterceptorFunc(func(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
					interceptorCalls = append(interceptorCalls, 3)
					return next(ctx, msgs...)
				}),
			}
			msg = "test message"
			bp  = buildInterceptors(p, interceptors)
		)

		// When
		err := bp.Produce(context.Background(), Message{
			Value: []byte(msg),
		})

		// Then
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, interceptorCalls)
		assert.Equal(t, msg, p.msg)
	})

	t.Run("it should build interceptors and returns producer error when interceptors provided", func(t *testing.T) {
		// Given
		var (
			p                = &testProducer{returnErr: true}
			interceptorCalls []int
			interceptors     = []Interceptor{
				InterceptorFunc(func(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
					interceptorCalls = append(interceptorCalls, 1)
					return next(ctx, msgs...)
				}),
				InterceptorFunc(func(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
					interceptorCalls = append(interceptorCalls, 2)
					return next(ctx, msgs...)
				}),
				InterceptorFunc(func(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
					interceptorCalls = append(interceptorCalls, 3)
					return next(ctx, msgs...)
				}),
			}
			msg = "test message"
			bp  = buildInterceptors(p, interceptors)
		)

		// When
		err := bp.Produce(context.Background(), Message{
			Value: []byte(msg),
		})

		// Then
		assert.Equal(t, errTest, err)
		assert.Equal(t, []int{1, 2, 3}, interceptorCalls)
		assert.Equal(t, "", p.msg)
	})
}
