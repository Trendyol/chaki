package newrelic

import (
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Module(providers ...any) *module.Module {
	return module.New("newrelic").Provide(
		newNr,
		providers,
	)
}

func newNr(cfg *config.Config) (*newrelic.Application, error) {
	of := cfg.Of("newrelic")

	// set default configs
	of.SetDefault("agentenabled", false)

	if !of.Exists("licensekey") {
		return newrelic.NewApplication(
			newrelic.ConfigAppLogEnabled(of.GetBool("agentenabled")),
			newrelic.ConfigEnabled(false),
		)
	}

	return newrelic.NewApplication(
		newrelic.ConfigLicense(of.GetString("licensekey")),
		newrelic.ConfigAppLogEnabled(of.GetBool("agentenabled")),
		newrelic.ConfigAppName(of.GetString("appname")),
	)
}
