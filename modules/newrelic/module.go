package newrelic

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ModuleName = "chaki-newrelic-module"

// config key consts
const (
	keyEnabled = "enabled"

	keyAppName = "appname"
	keyLicence = "licence"

	keyAppLogEnabled = "logenabled"
)

func Module(option ...Option) *module.Module {
	o := buildOptions(option...)
	return module.
		New("newrelic").
		Provide(
			chaki.Valuer(o),
			newNr,
		).
		Merge(
			o.subModules...,
		)
}

func newNr(cfg *config.Config, opts *options) (*newrelic.Application, error) {
	nrcfg := cfg.Of("newrelic")
	setDefaultConfigs(nrcfg)

	if !nrcfg.GetBool(keyEnabled) {
		return newrelic.NewApplication(append(opts.nrOptions, newrelic.ConfigEnabled(false))...)
	}

	opts.nrOptions = append(
		opts.nrOptions,
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName(nrcfg.GetString(keyAppName)),
		newrelic.ConfigLicense(nrcfg.GetString(keyAppName)),
		newrelic.ConfigAppLogEnabled(nrcfg.GetBool(keyAppLogEnabled)),
	)

	return newrelic.NewApplication(opts.nrOptions...)
}

func setDefaultConfigs(cfg *config.Config) {
	cfg.SetDefault(keyEnabled, false)
	cfg.SetDefault(keyLicence, "")
	cfg.SetDefault(keyAppLogEnabled, true)
}
