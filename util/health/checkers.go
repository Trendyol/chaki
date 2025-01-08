package health

type Checker interface {
	Check() error
}

type HookFunc func() error

func (hf HookFunc) Check() error {
	return hf()
}

type (
	LivenessChecker  Checker
	ReadinessChecker Checker
)

// type-check
var (
	_ LivenessChecker  = HookFunc(nil)
	_ ReadinessChecker = HookFunc(nil)
	_ Checker          = HookFunc(nil)
)
