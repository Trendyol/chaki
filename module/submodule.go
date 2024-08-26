package module

// SubModule represents a simple provide pool partition of a module
// can not be used as a standalone module with chaki.
// It has to be merged with another module with
type SubModule struct {
	provides []any
}

func NewSubModule() *SubModule {
	return &SubModule{
		provides: make([]any, 0),
	}
}

// Provide add constructor to provide pool
func (sm *SubModule) Provide(provide ...any) *SubModule {
	sm.provides = append(sm.provides, provide...)
	return sm
}
