package couchbase

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/couchbase/gocb/v2"
)

var (
	asClusterOptionsWrapper = as.Struct[ClusterOptionsWrapper]("cbclusteroptionswrapper")
	asCouchbaseTracers      = as.Interface[gocb.RequestTracer]("gocbrequesttracers")
)

func Module(option ...Option) *module.Module {
	m := module.New()

	m.Provide(
		chaki.Valuer(option),
		buildOptions,
		asClusterOptionsWrapper.Grouper(),
		asCouchbaseTracers.Grouper(),
		newCluster,
	)

	m.AddProvideHook(
		module.ProvideHook{
			Match: asClusterOptionsWrapper.Match,
			Wrap:  asClusterOptionsWrapper.Value,
		},
	)

	return m
}

func newCluster(cfg *config.Config, opts *options) (*gocb.Cluster, error) {
	cbcfg := cfg.Of("couchbase")
	setDefaultConfigs(cbcfg)

	clutserOptions := gocb.ClusterOptions{
		Username: cbcfg.GetString("username"),
		Password: cbcfg.GetString("password"),
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout:    cbcfg.GetDuration("connecttimeout"),
			KVTimeout:         cbcfg.GetDuration("kvtimeout"),
			KVDurableTimeout:  cbcfg.GetDuration("kvdurabletimeout"),
			KVScanTimeout:     cbcfg.GetDuration("kvscantimeout"),
			ViewTimeout:       cbcfg.GetDuration("viewtimeout"),
			QueryTimeout:      cbcfg.GetDuration("querytimeuot"),
			AnalyticsTimeout:  cbcfg.GetDuration("analyticstimeout"),
			SearchTimeout:     cbcfg.GetDuration("searchtimeout"),
			ManagementTimeout: cbcfg.GetDuration("managmenttimeuot"),
		},
	}

	for _, wr := range opts.clusterOptionsWrappers {
		clutserOptions = wr(clutserOptions)
	}

	if len(opts.tracers) > 0 {
		clutserOptions.Tracer = newJoinedTracer(opts.tracers)
	}

	return gocb.Connect(cbcfg.GetString("host"), clutserOptions)
}
