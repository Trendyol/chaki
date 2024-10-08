package couchbase

import (
	"context"
	"time"

	"github.com/Trendyol/chaki/util/slc"
	"github.com/couchbase/gocb/v2"
)

type emptySpanner struct{}

func (span emptySpanner) End()                                       {}
func (span emptySpanner) Context() gocb.RequestSpanContext           { return nil }
func (span emptySpanner) SetAttribute(key string, value interface{}) {}
func (span emptySpanner) AddEvent(key string, timestamp time.Time)   {}

type contextSpanner struct {
	emptySpanner
	ctx context.Context
}

func ParentSpan(ctx context.Context) gocb.RequestSpan {
	return &contextSpanner{
		ctx: ctx,
	}
}

func (span contextSpanner) Context() gocb.RequestSpanContext {
	return span.ctx
}

type joinedTracer struct {
	tracers []gocb.RequestTracer
}

func newJoinedTracer(tracers []gocb.RequestTracer) gocb.RequestTracer {
	return &joinedTracer{tracers}
}

func (jt *joinedTracer) RequestSpan(parentContext gocb.RequestSpanContext, operationName string) gocb.RequestSpan {
	spanners := slc.Map(jt.tracers, func(t gocb.RequestTracer) gocb.RequestSpan {
		return t.RequestSpan(parentContext, operationName)
	})

	return &joinedSpanner{spanners}
}

type joinedSpanner struct {
	spans []gocb.RequestSpan
}

func (s *joinedSpanner) End() {
	for _, span := range s.spans {
		span.End()
	}
}

func (s *joinedSpanner) Context() gocb.RequestSpanContext {
	if len(s.spans) > 0 {
		return s.spans[0].Context()
	}

	return context.Background()
}

func (s *joinedSpanner) AddEvent(name string, timestamp time.Time) {
	for _, span := range s.spans {
		span.AddEvent(name, timestamp)
	}
}

func (s *joinedSpanner) SetAttribute(key string, value interface{}) {
	for _, span := range s.spans {
		span.SetAttribute(key, value)
	}
}
