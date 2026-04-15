# Cluster Config And Tracing

Use this page when implementing or changing Couchbase-backed features.

## Module Pattern

- Register the module with `couchbase.Module(...)`.
- Inject `*gocb.Cluster` and the relevant config to resolve buckets and collections.
- Use module options such as `WithTracer(...)` and `WithClusterOptionWrapper(...)` when client behavior needs to be customized.

## Config Rules

- Keep host, bucket name, username, password, and timeout values in config.
- Reused bucket or collection names should live in constants rather than repeated string literals.

## Tracing And Health

- Use `couchbase.ParentSpan(ctx)` in request options where parent span support is available.
- Rely on the module-provided health probing behavior before adding duplicate connectivity checks in application code.
