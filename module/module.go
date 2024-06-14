package module

import (
	"strings"

	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/util/slc"
)

type condProvide struct {
	requireds string
	provides  []any
}

func (cp *condProvide) match(modules ...string) bool {
	rm := strings.Split(cp.requireds, ",")
	for _, rmi := range rm {
		if slc.Contains(modules, rmi) {
			return true
		}
	}
	return false
}

// ProvideHook represents a pre-Provide step so that constructors can be manipulated or grouped
type ProvideHook struct {
	Match func(ctr any) bool
	Wrap  func(ctr any) any
}

// Meta is the meta data of a module
type Meta struct {
	Name         string
	Provides     []any
	Invokes      []any
	ProvideHooks []ProvideHook
	Depends      []string
}

// Module represents a group of provide and invoke functions.
// It makes it possible to group functionalities with modules.
// These functionalities can be imported whenever needed
type Module struct {
	meta         Meta
	condProvides []condProvide
}

func New(name ...string) *Module {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return &Module{
		meta: Meta{
			Name: n,
		},
	}
}

// Provide add constructors to provide pool
func (m *Module) Provide(ctr ...any) *Module {
	m.meta.Provides = append(m.meta.Provides, ctr...)
	return m
}

// CondProvide add constructors to provide pool with module import conditions
// Example:
//
//	m.CondProvide("foo,bar", NewFooBarService) // provides if foo or bar module is used by application
//
// m.CondProvide("foo", NewFooController)
func (m *Module) CondProvide(requiredModules string, ctr ...any) {
	m.condProvides = append(m.condProvides, condProvide{requiredModules, ctr})
}

// Name returns name of the module
func (m *Module) Name() string {
	return m.meta.Name
}

// Invoke add functions to invoke pool
func (m *Module) Invoke(ctr ...any) *Module {
	m.meta.Invokes = append(m.meta.Invokes, ctr...)
	return m
}

// AddProvideHook adds provide hook
func (m *Module) AddProvideHook(hooks ...ProvideHook) *Module {
	m.meta.ProvideHooks = append(m.meta.ProvideHooks, hooks...)
	return m
}

// AddAsser adds provide hook from as.Asser instances
func (m *Module) AddAsser(assers ...as.Asser) *Module {
	for _, asser := range assers {
		m.Provide(asser.Grouper())
		m.AddProvideHook(ProvideHook{
			Match: asser.Match,
			Wrap:  asser.Value,
		})
	}
	return m
}

// Meta returns meta data of the module
func (m *Module) Meta(usedModules ...string) Meta {
	r := Meta{
		Name:         m.meta.Name,
		Provides:     m.meta.Provides,
		Invokes:      m.meta.Invokes,
		ProvideHooks: m.meta.ProvideHooks,
		Depends:      m.meta.Depends,
	}

	if len(usedModules) == 0 {
		return r
	}

	for _, v := range m.condProvides {
		if v.match(usedModules...) {
			r.Provides = append(r.Provides, v.provides...)
		}
	}

	return r
}
