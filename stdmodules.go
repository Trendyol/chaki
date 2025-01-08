package chaki

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/util/health"
	"go.uber.org/fx"
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

func loggerModule() *module.Module {
	syncLogger := func(lc fx.Lifecycle) {
		lc.Append(fx.StartStopHook(logger.Init, logger.Sync))
	}
	return module.New("logger").Invoke(syncLogger)
}

func healthModule() *module.Module {
	asLivenessChecker := as.Interface[health.LivenessChecker]("healthlivenesschecker")
	asReadinessChecker := as.Interface[health.ReadinessChecker]("healthreadinesschecker")

	m := module.New("health")

	m.Provide(
		asLivenessChecker.Grouper(),
		asReadinessChecker.Grouper(),
	)

	m.AddProvideHook(
		module.ProvideHook{
			Match: asLivenessChecker.Match,
			Wrap:  asLivenessChecker.Value,
		},
		module.ProvideHook{
			Match: asReadinessChecker.Match,
			Wrap:  asReadinessChecker.Value,
		},
	)

	return m
}
