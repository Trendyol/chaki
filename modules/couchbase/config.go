package couchbase

import (
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/wrapper"
	"github.com/couchbase/gocb/v2"
)

type ClusterOptionsWrapper = wrapper.Wrapper[gocb.ClusterOptions]

var defaultConfigMap = map[string]any{
	"host":     "",
	"username": "",
	"password": "",
	// timeout configs
	"connecttimeout":   "1000ms",
	"kvtimeout":        "2500ms",
	"kvdurabletimeout": "10000ms",
	"kvscantimeout":    "10000ms",
	"viewtimeout":      "75000ms",
	"querytimeuot":     "75000ms",
	"analyticstimeout": "75000ms",
	"searchtimeout":    "75000ms",
	"managmenttimeuot": "75000ms",
}

func setDefaultConfigs(cfg *config.Config) {
	for k, v := range defaultConfigMap {
		cfg.SetDefault(k, v)
	}
}
