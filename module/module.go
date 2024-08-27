package module

import (
	"github.com/Trendyol/chaki/as"
)

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
	meta Meta
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

// Merge combine sub module providers into module providers
func (m *Module) Merge(sub ...*SubModule) *Module {
	for _, s := range sub {
		m.Provide(s.provides...)
		m.AddProvideHook(s.provideHooks...)
		m.AddAsser(s.assers...)
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

	return r
}
