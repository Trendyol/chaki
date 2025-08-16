# Couchbase Module

The Couchbase module integrates Couchbase database support into Chaki applications, providing a configured *gocb.Cluster instance, health probes, and tracing capabilities.

## Usage

Add the module to your Chaki app, optionally with custom options:

```go
app.Use(couchbase.Module(
  couchbase.WithTracer(/* custom tracer */),
  couchbase.WithClusterOptionWrapper(/* wrapper func */),
))
```

The module provides a *gocb.Cluster for injection into repositories or services.

## Configuration

Configure the Couchbase connection and timeouts via YAML. Defaults are provided for timeouts.

Example (as provided):

```yaml
couchbase:
  host: "couchbase://127.0.0.1"
  bucketname: "Bucketname"
  connecttimeout: "5000ms"
  username: ${secret:couchbaseUsername}
  password: ${secret:couchbasePassword}
```

Full configuration options (with defaults):

- host: "" (required, e.g., "couchbase://localhost")
- username: "" (required)
- password: "" (required)
- connecttimeout: "1000ms"
- kvtimeout: "2500ms"
- kvdurabletimeout: "10000ms"
- kvscantimeout: "10000ms"
- viewtimeout: "75000ms"
- querytimeout: "75000ms"
- analyticstimeout: "75000ms"
- searchtimeout: "75000ms"
- managmenttimeout: "75000ms"

Secrets can be referenced as shown for username and password.

## Features

### Health Probes

Automatically provides liveness and readiness probes using cluster.Ping(). Implements health.Probe interface.

### Tracing

Supports multiple gocb.RequestTracer via WithTracer option. Use couchbase.ParentSpan(ctx) for request parent spans. Joined tracers are supported for combining multiple tracers.

### Cluster Options Wrappers

Customize gocb.ClusterOptions using WithClusterOptionWrapper to wrap and modify the options before connecting.

## Example

From the example/withcouchbase:

```go
// main.go
app := chaki.New()
app.Use(
  server.Module(),
  couchbase.Module(),
)
app.Provide(
  NewRepository,
  newController,
)
```

```go
// repository.go
type fooRepository struct {
  coll *gocb.Collection
}

func NewRepository(cluster *gocb.Cluster, cfg *config.Config) FooRepository {
  return &fooRepository{
    coll: cluster.Bucket(cfg.GetString("couchbase.bucketname")).DefaultCollection(),
  }
}

func (r *fooRepository) Get(ctx context.Context, id string) (*Foo, error) {
  resp, err := r.coll.Get(id, &gocb.GetOptions{
    ParentSpan: couchbase.ParentSpan(ctx), // Enable tracing
  })
  // ...
}
```

For more details, see module.go, config.go, option.go, and util.go.
