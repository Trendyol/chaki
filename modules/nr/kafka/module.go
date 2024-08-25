package nrkafka

import (
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ModuleName = "chaki-nr-kafka-module"

func Module() *module.Module {
	return module.New(ModuleName).Provide(newConsumerInterceptor, newBulkConsumerInterceptor)
}

func newConsumerInterceptor(nr *newrelic.Application) consumer.Interceptor {
	return consumer.InterceptorFunc(func(msg *consumer.Message, next consumer.ConsumeFn) error {
		txn := nr.StartTransaction(msg.Topic)
		defer txn.End()
		msg.Context = newrelic.NewContext(msg.Context, txn)

		if err := next(msg); err != nil {
			txn.NoticeError(err)
			return err
		}

		return nil
	})
}

func newBulkConsumerInterceptor(nr *newrelic.Application) consumer.BatchConsumerInterceptor {
	return consumer.BatchConsumerInterceptorFunc(func(msg []*consumer.Message, next consumer.BatchConsumeFn) error {
		if len(msg) == 0 {
			return next(msg)
		}

		txn := nr.StartTransaction(msg[0].Topic)
		defer txn.End()

		for i := range msg {
			msg[i].Context = newrelic.NewContext(msg[i].Context, txn)
		}

		if err := next(msg); err != nil {
			txn.NoticeError(err)
			return err
		}

		return nil
	})
}
