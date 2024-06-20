package consumer

import (
	"testing"

	"github.com/Trendyol/chaki/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_NewStarter(t *testing.T) {
	t.Run("it shouold build consumer with only consumer fn", func(t *testing.T) {
		// Given
		var (
			cfg    = config.NewConfig(viper.New(), nil)
			doneCh = make(chan bool, 1)
			c      = NewFn("example", func(m *Message) error {
				doneCh <- true
				return nil
			})
		)

		cfg.Set("kafka.consumer.example", map[string]any{
			"brokers": "text.kafka.com",
			"topic":   "topic-topic",
			"groupId": "gorup-id",
		})

		// When
		starter, err := NewStarter(cfg, c, nil)
		consumerErr := starter.consumerFn(nil)

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
			c         = NewFn("example", func(m *Message) error {
				resultArr = append(resultArr, 4)
				return nil
			})

			interceptors = []Interceptor{
				InterceptorFunc(func(msg *Message, next ConsumeFn) error {
					resultArr = append(resultArr, 1)
					return next(msg)
				}),
				InterceptorFunc(func(msg *Message, next ConsumeFn) error {
					resultArr = append(resultArr, 2)
					return next(msg)
				}),
				InterceptorFunc(func(msg *Message, next ConsumeFn) error {
					resultArr = append(resultArr, 3)
					return next(msg)
				}),
			}
		)

		cfg.Set("kafka.consumer.example", map[string]any{
			"brokers": "text.kafka.com",
			"topic":   "topic-topic",
			"groupId": "gorup-id",
		})

		// When
		starter, err := NewStarter(cfg, c, interceptors)
		consumerErr := starter.consumerFn(nil)

		// Then
		assert.NoError(t, err)
		assert.NoError(t, consumerErr)
		assert.Equal(t, []int{1, 2, 3, 4}, resultArr)
	})
}
