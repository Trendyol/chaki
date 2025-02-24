package couchbase

import (
	"context"
	"testing"
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/couchbase/gocb/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewLivenessReadinessProbe(t *testing.T) {
	cl := &gocb.Cluster{}
	patch := gomonkey.ApplyMethod(cl, "Ping", func(_ *gocb.Cluster, opts *gocb.PingOptions) (*gocb.PingResult, error) {
		return &gocb.PingResult{}, nil
	})
	defer patch.Reset()

	probe := newLivenessReadinessProbe(cl)

	errLiveness := probe.Liveness(context.Background())
	errReadiness := probe.Readiness(context.Background())

	assert.NoError(t, errLiveness)
	assert.NoError(t, errReadiness)
}

func Test_newCluster(t *testing.T) {
	// given
	v := viper.New()
	cfg := config.NewConfig(v, nil)

	cfg.Set("couchbase.host", "foo_bar_host")
	cfg.Set("couchbase.username", "foo")
	cfg.Set("couchbase.password", "bar")
	cfg.Set("couchbase.connecttimeout", "2ms")
	cfg.Set("couchbase.kvtimeout", "3ms")
	cfg.Set("couchbase.querytimeout", "4ms")
	cfg.Set("couchbase.searchtimeout", "5ms")

	called := false

	patch := gomonkey.ApplyFunc(gocb.Connect, func(c string, opts gocb.ClusterOptions) (*gocb.Cluster, error) {
		called = true
		assert.Equal(t, "foo_bar_host", c)
		assert.Equal(t, "foo", opts.Username)
		assert.Equal(t, "bar", opts.Password)
		assert.Equal(t, 2*time.Millisecond, opts.TimeoutsConfig.ConnectTimeout)
		return nil, nil
	})
	defer patch.Reset()

	// when
	_, err := newCluster(cfg, buildOptions(nil, nil, nil))

	// then
	assert.NoError(t, err)
	assert.True(t, called)
}
