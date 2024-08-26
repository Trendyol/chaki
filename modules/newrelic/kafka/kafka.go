package kafka

import (
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	nrmodule "github.com/Trendyol/chaki/modules/newrelic"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func WithKafka() nrmodule.Option {
	sm := module.NewSubModule().Provide(newKafkaInterceptor)
	return nrmodule.WithSubModule(sm)
}

type Interceptor struct {
	nr *newrelic.Application
}

func newKafkaInterceptor(nr *newrelic.Application) consumer.Interceptor {
	return &Interceptor{nr: nr}
}

func (k *Interceptor) Intercept(msg *consumer.Message, next consumer.ConsumeFn) error {
	txn := k.nr.StartTransaction(msg.Topic)
	defer txn.End()

	txn.AddAttribute("kafka.partition", msg.Partition)
	txn.AddAttribute("kafka.offset", msg.Offset)
	txn.AddAttribute("kafka.key", string(msg.Key))
	txn.AddAttribute("kafka.headers", msg.Headers)
	txn.AddAttribute("kafka.time", msg.Time)

	ctx := newrelic.NewContext(msg.Context, txn)
	msg.Context = ctx

	err := next(msg)
	if err != nil {
		txn.NoticeError(err)
		txn.AddAttribute("error.message", err.Error())
	}

	return err
}
