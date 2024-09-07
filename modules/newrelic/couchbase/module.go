package couchbase

import (
	"context"
	"time"

	"github.com/Trendyol/chaki/module"
	nrmodule "github.com/Trendyol/chaki/modules/newrelic"
	"github.com/couchbase/gocb/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func WithCouchbase() nrmodule.Option {
	return nrmodule.WithSubModule(module.NewSubModule().Provide(newCbTracer))
}

type cbTracer struct{}

func newCbTracer() gocb.RequestTracer {
	return &cbTracer{}
}

func noopSpan() gocb.RequestSpan {
	return new(gocb.NoopTracer).RequestSpan(context.Background(), "")
}

func (cbTracer) RequestSpan(parentContext gocb.RequestSpanContext, operationName string) gocb.RequestSpan {
	if parentContext == nil {
		return noopSpan()
	}

	parent := context.Background()
	if ctx, ok := parentContext.(context.Context); ok {
		parent = ctx
	}

	txn := newrelic.FromContext(parent)
	if txn == nil {
		return noopSpan()
	}

	sgm := &newrelic.DatastoreSegment{
		StartTime: txn.StartSegmentNow(),
		Product:   newrelic.DatastoreCouchDB,
		Operation: operationName,
	}

	return newCBSpan(sgm, parent)
}

type cbSpan struct {
	sgm *newrelic.DatastoreSegment
	ctx context.Context
}

func newCBSpan(sgm *newrelic.DatastoreSegment, ctx context.Context) gocb.RequestSpan {
	return &cbSpan{
		sgm: sgm,
		ctx: ctx,
	}
}

func (s *cbSpan) End() {
	s.sgm.End()
}

func (s *cbSpan) Context() gocb.RequestSpanContext {
	return s.ctx
}

func (s *cbSpan) AddEvent(name string, timestamp time.Time) {
	newrelic.FromContext(s.ctx).RecordLog(newrelic.LogData{
		Message:   name,
		Timestamp: timestamp.UnixMilli(),
	})
}

func (s *cbSpan) SetAttribute(key string, value interface{}) {
	s.sgm.AddAttribute(key, value)
}
