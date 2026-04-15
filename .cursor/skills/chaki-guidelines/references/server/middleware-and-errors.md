# Middleware And Errors

Use this page when the task affects request middlewares, logging enrichment, or server error shaping.

## Middleware Rules

- Use controller-level `AddMiddleware(...)` or route-level `AddMiddlewares(...)` when attaching request middleware.
- Middleware that enriches logger or correlation data should start from `c.UserContext()`, add fields to that context, and write it back with `c.SetUserContext(ctx)` before continuing.
- Keep middleware focused on boundary concerns such as logging, tracing, auth, correlation, and request metadata.

## Error Handling Rules

- Register custom server error handlers through `server.WithErrorHandler(...)`.
- Include request metadata such as correlation ID, request URI, and request method when building structured transport errors.
- Use the existing context-based logger for error logging so request metadata stays attached.
- Log expected user-facing failures at the right level instead of treating every error as an unexpected internal fault.

## Practical Guidance

- If a response type needs dynamic HTTP status selection, implement `Status() int` on the response object rather than scattering manual status writes.
- Keep raw internal error strings out of public responses unless that behavior is explicitly intended.
- When error-related header names or context keys are reused, extract them into constants.
