package chaki

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/util/health"
)

func configModule(opts configOptions) (*module.Module, error) {
	m := module.New("config")
	if opts.disabled {
		return m, nil
	}

	cfg, err := config.NewConfigFromPaths(opts.path, opts.referencePaths)
	if err != nil {
		return nil, err
	}

	m.Provide(Valuer(cfg))

	return m, nil
}

func healthModule() *module.Module {
	asHealthProbes := as.Interface[health.Probe]("healthprobes")

	m := module.New("health")

	m.Provide(
		asHealthProbes.Grouper(),
	)

	m.AddProvideHook(
		module.ProvideHook{
			Match: asHealthProbes.Match,
			Wrap:  asHealthProbes.Value,
		},
	)

	return m
}
