package couchbase

import "github.com/couchbase/gocb/v2"

type options struct {
	clusterOptionsWrappers []ClusterOptionsWrapper
	tracers                []gocb.RequestTracer
}

type Option interface {
	Apply(*options)
}

type withOption func(*options)

func (wf withOption) Apply(opts *options) {
	wf(opts)
}

func WithTracer(tr ...gocb.RequestTracer) Option {
	return withOption(func(o *options) {
		o.tracers = append(o.tracers, tr...)
	})
}

func WithClusterOptionWrapper(wrappers ...ClusterOptionsWrapper) Option {
	return withOption(func(o *options) {
		o.clusterOptionsWrappers = append(o.clusterOptionsWrappers, wrappers...)
	})
}

func buildOptions(opt []Option, clusterOptionWrappers []ClusterOptionsWrapper, tracers []gocb.RequestTracer) *options {
	opts := &options{}
	for _, o := range opt {
		o.Apply(opts)
	}

	WithClusterOptionWrapper(clusterOptionWrappers...).Apply(opts)
	WithTracer(tracers...).Apply(opts)

	return opts
}
