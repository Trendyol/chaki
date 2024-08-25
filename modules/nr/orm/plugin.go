package nrorm

import (
	"github.com/Trendyol/chaki/util/appctx"
	"github.com/Trendyol/chaki/util/wrapper"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

type SegmentHook wrapper.Wrapper[*newrelic.DatastoreSegment]

type gormHookFunc func(tx *gorm.DB)

type gormRegister interface {
	Register(name string, fn func(*gorm.DB)) error
}

var segmentCtxValuer = appctx.NewValuer[*newrelic.DatastoreSegment]("gormnewrelicsegment")

// gormPlugin is an gorm.Plugin implementtation for newrelicc support of the chaki framework
// this plugin is newrelic version of the gorm's official otel plugin
// See https://github.com/go-gorm/opentelemetry/blob/master/tracing/tracing.go
type gormPlugin struct {
	nr    *newrelic.Application
	hooks []SegmentHook
}

func newGormPlugin(nr *newrelic.Application, hooks []SegmentHook) gorm.Plugin {
	return &gormPlugin{nr, hooks}
}

func (p *gormPlugin) Name() string {
	return "chaki-newrelic-gorm"
}

func (p *gormPlugin) Initialize(db *gorm.DB) error {
	cb := db.Callback()
	hooks := []struct {
		callback gormRegister
		hook     gormHookFunc
		name     string
	}{
		{cb.Create().Before("gorm:create"), p.before("gorm.Create"), "before:create"},
		{cb.Create().After("gorm:create"), p.after(), "after:create"},

		{cb.Query().Before("gorm:query"), p.before("gorm.Query"), "before:select"},
		{cb.Query().After("gorm:query"), p.after(), "after:select"},

		{cb.Delete().Before("gorm:delete"), p.before("gorm.Delete"), "before:delete"},
		{cb.Delete().After("gorm:delete"), p.after(), "after:delete"},

		{cb.Update().Before("gorm:update"), p.before("gorm.Update"), "before:update"},
		{cb.Update().After("gorm:update"), p.after(), "after:update"},

		{cb.Row().Before("gorm:row"), p.before("gorm.Row"), "before:row"},
		{cb.Row().After("gorm:row"), p.after(), "after:row"},

		{cb.Raw().Before("gorm:raw"), p.before("gorm.Raw"), "before:raw"},
		{cb.Raw().After("gorm:raw"), p.after(), "after:raw"},
	}

	for _, h := range hooks {
		if err := h.callback.Register("nr:"+h.name, h.hook); err != nil {
			return err
		}
	}

	return nil
}

func (p *gormPlugin) before(_ string) gormHookFunc {
	return func(tx *gorm.DB) {
		seg := &newrelic.DatastoreSegment{
			StartTime: newrelic.FromContext(tx.Statement.Context).StartSegmentNow(),
		}
		tx.Statement.Context = segmentCtxValuer.Set(tx.Statement.Context, seg)
	}
}

func (p *gormPlugin) after() gormHookFunc {
	return func(tx *gorm.DB) {
		seg := p.buildSegment(tx)
		for _, h := range p.hooks {
			seg = h(seg)
		}

		defer seg.End()

	}
}

func (p *gormPlugin) buildSegment(tx *gorm.DB) *newrelic.DatastoreSegment {
	seg := segmentCtxValuer.Get(tx.Statement.Context)
	if seg == nil {
		return nil
	}

	// Add attributes
	seg.Collection = tx.Statement.Table
	seg.RawQuery = tx.Statement.SQL.String()
	seg.Product = newrelic.DatastoreProduct(tx.Dialector.Name())
	seg.DatabaseName = tx.Statement.DB.Name()

	return seg
}
