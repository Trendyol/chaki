package couchbase

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/couchbase/gocb/v2"
)

var (
	asClusterOptionsWrapper = as.Struct[ClusterOptionsWrapper]("cbclusteroptionswrapper")
	asCouchbaseTracers      = as.Interface[gocb.RequestTracer]("gocbrequesttracers")
)

func Module() *module.Module {
	m := module.New()

	m.Provide(
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

func newCluster(cfg *config.Config, wrappers []ClusterOptionsWrapper, tracers []gocb.RequestTracer) (*gocb.Cluster, error) {
	cbcfg := cfg.Of("couchbase")
	setDefaultConfigs(cbcfg)

	opts := gocb.ClusterOptions{
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

	for _, wr := range wrappers {
		opts = wr(opts)
	}

	if len(tracers) > 0 {
		opts.Tracer = newJoinedTracer(tracers)
	}

	return gocb.Connect(cbcfg.GetString("host"), opts)
}
