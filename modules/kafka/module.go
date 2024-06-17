package kafka

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

const ModuleName = "chaki-kafka"

var (
	asConsumer = as.Interface[consumer.Consumer]("kafkaconsumers")
)

func Module() *module.Module {
	m := module.New(ModuleName)

	m.Invoke(runConsumers)

	m.AddAsser(asConsumer)

	return m
}

func runConsumers(consumers []consumer.Consumer, cfg *config.Config, lc fx.Lifecycle, eg *errgroup.Group) error {
	starters := make([]*consumer.Starter, len(consumers))
	for i, c := range consumers {
		s, err := consumer.NewStarter(cfg, c)
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
