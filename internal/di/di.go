package di

import (
	"time"

	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Injector struct {
	provides  []any
	decorates []any
	invokes   []any
}

func NewInjector() *Injector {
	return &Injector{}
}

func (d *Injector) Provide(ctr ...any) {
	d.provides = append(d.provides, ctr...)
}

func (d *Injector) Decorate(ctr ...any) {
	d.decorates = append(d.decorates, ctr...)
}

func (d *Injector) Invoke(ctr ...any) {
	d.invokes = append(d.invokes, ctr...)
}

func (d *Injector) Start(timeout time.Duration) error {
	eg := &errgroup.Group{}
	app := fx.New(
		fx.Provide(ByValue(eg)),
		fx.Provide(d.provides...),
		fx.Decorate(d.decorates...),
		fx.Invoke(d.invokes...),
		fx.StartTimeout(timeout),
	)

	go func() {
		if err := eg.Wait(); err != nil {
			panic(err)
		}
	}()

	app.Run()

	return app.Err()
}
