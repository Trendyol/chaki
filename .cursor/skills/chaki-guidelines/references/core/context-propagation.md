# Context Propagation

Use this page when work crosses HTTP handlers, clients, repositories, Kafka, Couchbase, logging, or tracing.

## Main Rule

Preserve the incoming `context.Context` all the way through the call chain. Do not replace it with `context.Background()` or a fresh context in normal request flow.

## By Module

- Server typed handlers already receive `ctx context.Context`; pass that context into services.
- Fiber middlewares should read `c.UserContext()`, enrich it, then write it back with `c.SetUserContext(ctx)`.
- Client requests should use `Base.Request(ctx)` or `RequestWithCommand(ctx, commandName)` so outbound calls keep metadata and tracing.
- ORM access should go through `GormProvider.Get(ctx)` and transactional work should run through an injected `tx.Transactioner`, for example `transactioner.Transaction(parentCtx, func(ctx context.Context) error { ... })`.
- Kafka producers should produce with the caller context. Kafka consumers should preserve `msg.Context` for logs and downstream calls.
- Couchbase operations should attach `couchbase.ParentSpan(ctx)` where request options support parent spans.
- Logging should use `logger.From(ctx)` so contextual fields stay attached.

## Practical Guidance

- Add `ctx context.Context` as the first parameter of service, repository, and client methods.
- If a middleware enriches logging or correlation data, update the user context before calling the next handler.
- When names for headers, commands, or module keys are reused, define constants instead of repeating string literals.
