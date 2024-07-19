package consumer

import (
	"testing"

	"github.com/Trendyol/chaki/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewBatchStarter(t *testing.T) {
	t.Run("it shouold build consumer with only consumer fn", func(t *testing.T) {
		// Given
		var (
			cfg    = config.NewConfig(viper.New(), nil)
			doneCh = make(chan bool, 1)
			c      = NewBatchConsumerFn("example", func(m Messages) error {
				doneCh <- true
				return nil
			})
		)

		cfg.Set("kafka.consumer.example", map[string]any{
			"brokers": "text.kafka.com",
			"topic":   "topic-topic",
			"groupId": "gorup-id",
			"batchConfiguration": map[string]any{
				"messageGroupLimit": "2",
			},
		})

		// When
		starter, err := NewBatchConsumerStarter(cfg, c, nil)
		batchStarter, ok := starter.(*BatchConsumerStarter)
		assert.True(t, ok)

		consumerErr := batchStarter.consumerFn(nil)

		// Then
		assert.NoError(t, err)
		assert.NoError(t, consumerErr)
		<-doneCh
	})

	t.Run("it shouold build consumer with interceptor", func(t *testing.T) {
		// Given
		var (
			cfg       = config.NewConfig(viper.New(), nil)
			resultArr []int
			c         = NewBatchConsumerFn("example", func(m Messages) error {
				resultArr = append(resultArr, 4)
				return nil
			})

			interceptors = []BatchConsumerInterceptor{
				BatchConsumerInterceptorFunc(func(msg []*Message, next BatchConsumeFn) error {
					resultArr = append(resultArr, 1)
					return next(msg)
				}),
				BatchConsumerInterceptorFunc(func(msg []*Message, next BatchConsumeFn) error {
					resultArr = append(resultArr, 2)
					return next(msg)
				}),
				BatchConsumerInterceptorFunc(func(msg []*Message, next BatchConsumeFn) error {
					resultArr = append(resultArr, 3)
					return next(msg)
				}),
			}
		)

		cfg.Set("kafka.consumer.example", map[string]any{
			"brokers": "text.kafka.com",
			"topic":   "topic-topic",
			"groupId": "gorup-id",
			"batchConfiguration": map[string]any{
				"messageGroupLimit": "2",
			},
		})

		// When
		starter, err := NewBatchConsumerStarter(cfg, c, interceptors)

		batchStarter, ok := starter.(*BatchConsumerStarter)
		assert.True(t, ok)

		consumerErr := batchStarter.consumerFn(nil)

		// Then
		assert.NoError(t, err)
		assert.NoError(t, consumerErr)
		assert.Equal(t, []int{1, 2, 3, 4}, resultArr)
	})
}
