package newrelic

import "github.com/Trendyol/chaki/module"

type options struct {
	subModules []*module.SubModule
}

type Option interface {
	Apply(*options)
}

type withOption func(*options)

func (wf withOption) Apply(opts *options) {
	wf(opts)
}

func WithSubModule(m *module.SubModule) Option {
	return withOption(func(o *options) {
		o.subModules = append(o.subModules, m)
	})
}

func buildOptions(option ...Option) *options {
	opts := &options{}

	for _, o := range option {
		o.Apply(opts)
	}

	return opts
}
