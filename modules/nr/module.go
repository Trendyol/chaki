package nr

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
	keyLicence = "license"

	keyAppLogEnabled = "logenabled"
)

func Module(
	opts ...newrelic.ConfigOption,
) *module.Module {
	m := module.New(ModuleName)

	m.Provide(
		newApplication,
		chaki.Valuer(opts),
	)

	return m
}

func newApplication(cfg *config.Config, options []newrelic.ConfigOption) (*newrelic.Application, error) {
	nrcfg := cfg.Of("newrelic")
	setDefaultConfigs(nrcfg)

	if !nrcfg.GetBool(keyEnabled) {
		return newrelic.NewApplication(append(options, newrelic.ConfigEnabled(false))...)
	}

	o := append(
		options,
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName(nrcfg.GetString(keyAppName)),
		newrelic.ConfigLicense(nrcfg.GetString(keyAppName)),
		newrelic.ConfigAppLogEnabled(nrcfg.GetBool(keyAppLogEnabled)),
	)

	return newrelic.NewApplication(o...)

}

func setDefaultConfigs(cfg *config.Config) {
	cfg.SetDefault(keyEnabled, false)
	cfg.SetDefault(keyLicence, "")
	cfg.SetDefault(keyAppLogEnabled, true)
}
