package consumer

import (
	"github.com/Trendyol/chaki/config"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type Starter struct {
	libc kafkalib.Consumer
}

func NewStarter(cfg *config.Config, c Consumer) (*Starter, error) {
	libcfg, err := buildLibConfig(cfg, c)
	if err != nil {
		return nil, err
	}

	libcfg.ConsumeFn = c.Consume

	libc, err := kafkalib.NewConsumer(libcfg)
	if err != nil {
		return nil, err
	}

	return &Starter{libc}, nil
}

func (s *Starter) Start() error {
	s.libc.Consume()
	return nil
}

func (s *Starter) Stop() error {
	return s.libc.Stop()
}
