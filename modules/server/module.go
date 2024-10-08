package server

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/server/common"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

var (
	asController         = as.Interface[controller.Controller]("servercontrollers")
	asMiddleware         = as.Struct[fiber.Handler]("servermiddleware")
	asValidationRule     = as.Interface[validation.Rule]("validationrules")
	asFiberAppWrapper    = as.Struct[common.FiberAppWrapper]("fiberappwrappers")
	asMiddlewareGroup    = as.Struct[common.MiddlewareGroup]("middlewaregroups")
	asFiberConfigWrapper = as.Struct[common.FiberConfigWrapper]("fiberconfigwrappers")
)

func Module(option ...Option) *module.Module {
	opts := buildOptions(option...)

	m := module.New("server")

	m.Provide(
		// provide options
		chaki.Valuer(opts),

		// create registries
		asController.Handler(parseControllers),

		// swagger defs
		getSwaggerDefs,

		// middlewares - group
		asMiddleware.Grouper(),

		// validation - rules
		asValidationRule.Grouper(),

		asFiberAppWrapper.Grouper(),
		asMiddlewareGroup.Grouper(),
		asFiberConfigWrapper.Grouper(),

		// server
		defaultFiber,
		New,
	)

	m.Invoke(
		validation.Init,
		startServer,
	)

	m.AddProvideHook(
		module.ProvideHook{
			Match: asController.Match,
			Wrap:  asController.Value,
		},
		module.ProvideHook{
			Match: asMiddleware.Match,
			Wrap:  asMiddleware.Value,
		},
		module.ProvideHook{
			Match: asValidationRule.Match,
			Wrap:  asValidationRule.Value,
		},
		module.ProvideHook{
			Match: asFiberAppWrapper.Match,
			Wrap:  asFiberAppWrapper.Value,
		},
		module.ProvideHook{
			Match: asMiddlewareGroup.Match,
			Wrap:  asMiddlewareGroup.Value,
		},
		module.ProvideHook{
			Match: asFiberConfigWrapper.Match,
			Wrap:  asFiberConfigWrapper.Value,
		},
	)

	return m
}

func startServer(lc fx.Lifecycle, sv *Server, eg *errgroup.Group) {
	lc.Append(fx.StopHook(sv.Shutdown))
	eg.Go(func() error {
		return sv.Start()
	})
}
