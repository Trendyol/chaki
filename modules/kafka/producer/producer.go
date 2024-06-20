package producer

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/store"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

const (
	producerKey = "kafka.producer"
)

type (
	Message = kafkalib.Message

	Producer interface {
		Produce(context.Context, ...Message) error
	}

	producerFunc func(ctx context.Context, m ...Message) error

	InterceptNextFunc func(ctx context.Context, msg ...Message) error

	Interceptor interface {
		Intercept(ctx context.Context, msgs []Message, next InterceptNextFunc) error
	}

	InterceptorFunc func(ctx context.Context, msgs []Message, next InterceptNextFunc) error

	Factory interface {
		Get(name string) Producer
	}
)

var (
	_ Interceptor = InterceptorFunc(nil)
	_ Producer    = producerFunc(nil)
)

func (f producerFunc) Produce(ctx context.Context, msg ...Message) error {
	return f(ctx, msg...)
}

func (f InterceptorFunc) Intercept(ctx context.Context, msgs []Message, next InterceptNextFunc) error {
	return f(ctx, msgs, next)
}

type factory struct {
	producers *store.Bucket[string, Producer]
}

func NewFactory(cfg *config.Config, interceptors []Interceptor) (Factory, error) {
	b := store.NewBucket(func(k string) Producer { return nil })

	if !cfg.Exists(producerKey) {
		return &factory{b}, nil
	}

	pm := cfg.GetStringMap(producerKey)
	for name := range pm {
		p, err := newProducer(cfg, name, interceptors)
		if err != nil {
			return nil, err
		}
		b.Set(name, p)
	}

	return &factory{b}, nil
}

func (f *factory) Get(name string) Producer {
	return f.producers.Get(name)
}

func newProducer(cfg *config.Config, name string, interceptors []Interceptor) (Producer, error) {
	libcfg, err := buildLibConfig(cfg, name)
	if err != nil {
		return nil, err
	}

	lipp, err := kafkalib.NewProducer(libcfg)
	if err != nil {
		return nil, err
	}

	p := producerFunc(func(ctx context.Context, msg ...Message) error {
		return lipp.ProduceBatch(ctx, msg)
	})

	return buildInterceptors(p, interceptors), nil
}

func buildInterceptors(p Producer, interceptors []Interceptor) Producer {
	if len(interceptors) == 0 {
		return p
	}

	next := buildNextFunc(p, interceptors, 0)
	return producerFunc(next)
}

func buildNextFunc(p Producer, interceptors []Interceptor, i int) InterceptNextFunc {
	if i > len(interceptors)-1 {
		return func(ctx context.Context, msg ...Message) error {
			return p.Produce(ctx, msg...)
		}
	}

	next := buildNextFunc(p, interceptors, i+1)
	return func(ctx context.Context, msg ...Message) error {
		return interceptors[i].Intercept(ctx, msg, next)
	}
}
