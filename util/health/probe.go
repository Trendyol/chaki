package health

import "context"

type Probe interface {
	Liveness(ctx context.Context) error
	Readiness(ctx context.Context) error
}

type probe struct {
	liveness  func(ctx context.Context) error
	readiness func(ctx context.Context) error
}

func (p *probe) Liveness(ctx context.Context) error {
	if p.liveness == nil {
		return nil
	}
	return p.liveness(ctx)
}

func (p *probe) Readiness(ctx context.Context) error {
	if p.readiness == nil {
		return nil
	}
	return p.readiness(ctx)
}

func NewProbe(liveness, readiness func(ctx context.Context) error) Probe {
	return &probe{
		liveness:  liveness,
		readiness: readiness,
	}
}
