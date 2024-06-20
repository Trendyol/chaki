package consumer

import (
	"github.com/Trendyol/chaki/config"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type Starter struct {
	consumerFn ConsumeFn
	libc       kafkalib.Consumer
}

func NewStarter(cfg *config.Config, c Consumer, interceptors []Interceptor) (*Starter, error) {
	libcfg, err := buildLibConfig(cfg, c)
	if err != nil {
		return nil, err
	}

	libcfg.ConsumeFn = buildConsumeFn(c.Consume, interceptors)

	libc, err := kafkalib.NewConsumer(libcfg)
	if err != nil {
		return nil, err
	}

	return &Starter{libcfg.ConsumeFn, libc}, nil
}

func (s *Starter) Start() error {
	s.libc.Consume()
	return nil
}

func (s *Starter) Stop() error {
	return s.libc.Stop()
}

func buildConsumeFn(cf ConsumeFn, interceptors []Interceptor) ConsumeFn {
	if len(interceptors) == 0 {
		return cf
	}

	next := buildNextFunc(cf, interceptors, 1)
	return func(m *kafkalib.Message) error {
		return interceptors[0].Intercept(m, next)
	}

}

func buildNextFunc(cf ConsumeFn, interceptors []Interceptor, i int) ConsumeFn {
	if i > len(interceptors)-1 {
		return func(m *kafkalib.Message) error { return cf(m) }
	}

	next := buildNextFunc(cf, interceptors, i+1)
	return func(m *kafkalib.Message) error {
		return interceptors[i].Intercept(m, next)
	}
}
