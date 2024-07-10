package chaki

import (
	"time"
)

type configOptions struct {
	disabled       bool
	path           string
	referencePaths map[string]string
}

type options struct {
	timeout       time.Duration
	configOptions configOptions
}

func getOptions(opt ...Option) *options {
	def := &options{
		timeout: 45 * time.Second,
		configOptions: configOptions{
			path:           "resources/configs/application.yaml",
			referencePaths: map[string]string{},
		},
	}

	for _, o := range opt {
		o.Apply(def)
	}

	return def
}

type Option interface {
	Apply(*options)
}

type withOption func(*options)

func (wo withOption) Apply(opts *options) {
	wo(opts)
}

func WithTimeout(t time.Duration) Option {
	return withOption(func(o *options) {
		o.timeout = t
	})
}

func WithConfigDisabled(disabled bool) Option {
	return withOption(func(o *options) {
		o.configOptions.disabled = disabled
	})
}

func WithConfigPath(path string) Option {
	return withOption(func(o *options) {
		o.configOptions.path = path
	})
}

func WithConfigReferencePath(key, path string) Option {
	return withOption(func(o *options) {
		o.configOptions.referencePaths[key] = path
	})
}
