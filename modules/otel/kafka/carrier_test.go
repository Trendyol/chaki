package otelkafka

import (
	"testing"

	"github.com/Trendyol/chaki/modules/kafka/consumer"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
	"github.com/stretchr/testify/assert"
)

func Test_textMapCarrier_Get(t *testing.T) {
	// Given
	var (
		key    = "test-key"
		value  = "test-value"
		header = kafkalib.Header{Key: key, Value: []byte(value)}
		msg    = &consumer.Message{Headers: []kafkalib.Header{header}}
		mc     = NewMessageCarrier(msg)
	)

	// When
	got := mc.Get(key)

	// Then
	assert.Equal(t, value, got)
}

func Test_textMapCarrier_Set(t *testing.T) {
	// Given
	var (
		key   = "test-key"
		value = "test-value"
		msg   = &consumer.Message{Headers: []kafkalib.Header{}}
		mc    = NewMessageCarrier(msg)
	)

	// When
	mc.Set(key, value)

	// Then
	assert.Equal(t, value, string(msg.Headers[0].Value))
}

func Test_textMapCarrier_Keys(t *testing.T) {
	// Given
	var (
		msg = &consumer.Message{Headers: []kafkalib.Header{
			{Key: "key-1", Value: nil},
			{Key: "key-2", Value: nil},
			{Key: "key-3", Value: nil},
		}}
		mc = NewMessageCarrier(msg)
	)

	mc.Set("key-4", "")

	// When
	keys := mc.Keys()

	// Then
	assert.Equal(t, []string{"key-1", "key-2", "key-3", "key-4"}, keys)
}
