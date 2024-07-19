package kafka

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/kafka/producer"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

const ModuleName = "chaki-kafka"

var (
	asConsumer                 = as.Interface[consumer.Consumer]("kafkaconsumers")
	asBatchConsumer            = as.Interface[consumer.BatchConsumer]("kafkabatchconsumers")
	asConsumerInterceptor      = as.Interface[consumer.Interceptor]("consumerinterceptors")
	asBatchConsumerInterceptor = as.Interface[consumer.BatchConsumerInterceptor]("batchconsumerinterceptors")
	asProducerInterceptor      = as.Interface[producer.Interceptor]("producerinterceptors")
)

func Module() *module.Module {
	m := module.New(ModuleName)

	m.Provide(producer.NewFactory)

	m.Invoke(runConsumers)

	m.AddAsser(asConsumer, asBatchConsumer, asConsumerInterceptor, asBatchConsumerInterceptor, asProducerInterceptor)

	return m
}

func runConsumers(
	consumers []consumer.Consumer,
	batchConsumers []consumer.BatchConsumer,
	interceptors []consumer.Interceptor,
	batchConsumerInterceptors []consumer.BatchConsumerInterceptor,
	cfg *config.Config,
	lc fx.Lifecycle,
	eg *errgroup.Group,
) error {
	starters := make([]consumer.Starter, len(consumers)+len(batchConsumers))
	for i, c := range consumers {
		s, err := consumer.NewStarter(cfg, c, interceptors)
		if err != nil {
			return err
		}
		starters[i] = s
	}

	for i, bc := range batchConsumers {
		s, err := consumer.NewBatchConsumerStarter(cfg, bc, batchConsumerInterceptors)
		if err != nil {
			return err
		}
		starters[i] = s
	}

	for _, s := range starters {
		eg.Go(s.Start)
		lc.Append(fx.StopHook(s.Stop))
	}

	return nil
}
