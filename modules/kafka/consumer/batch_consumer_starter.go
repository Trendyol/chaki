package consumer

import (
	"errors"

	"github.com/Trendyol/chaki/config"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type BatchConsumerStarter struct {
	consumerFn BatchConsumeFn
	libc       kafkalib.Consumer
}

func NewBatchConsumerStarter(cfg *config.Config, c BatchConsumer, interceptors []BatchConsumerInterceptor) (Starter, error) {
	libcfg, err := buildLibConfig(cfg, c.Name())
	if err != nil {
		return nil, err
	}

	if libcfg.BatchConfiguration == nil {
		return nil, errors.New("batch consumer not configured") // todo add error types
	}

	libcfg.BatchConfiguration.BatchConsumeFn = buildBatchConsumeFn(c.Consume, interceptors)

	libc, err := kafkalib.NewConsumer(libcfg)
	if err != nil {
		return nil, err
	}

	return &BatchConsumerStarter{libcfg.BatchConfiguration.BatchConsumeFn, libc}, nil
}

func (s *BatchConsumerStarter) Start() error {
	s.libc.Consume()
	return nil
}

func (s *BatchConsumerStarter) Stop() error {
	return s.libc.Stop()
}

func buildBatchConsumeFn(cf BatchConsumeFn, interceptors []BatchConsumerInterceptor) BatchConsumeFn {
	if len(interceptors) == 0 {
		return cf
	}

	next := buildBatchConsumeNextFunc(cf, interceptors, 1)
	return func(m []*kafkalib.Message) error {
		return interceptors[0].Intercept(m, next)
	}
}

func buildBatchConsumeNextFunc(cf BatchConsumeFn, interceptors []BatchConsumerInterceptor, i int) BatchConsumeFn {
	if i > len(interceptors)-1 {
		return func(m []*kafkalib.Message) error { return cf(m) }
	}

	next := buildBatchConsumeNextFunc(cf, interceptors, i+1)
	return func(m []*kafkalib.Message) error {
		return interceptors[i].Intercept(m, next)
	}
}
