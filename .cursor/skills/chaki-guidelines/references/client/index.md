# Client Module

Use this page when the task adds or changes outbound HTTP integrations.

## Read Next

- Factory usage and client config shape: [factory-and-config.md](factory-and-config.md)
- Wrappers, retries, circuit breaker, and error decoding: [wrappers-retry-and-errors.md](wrappers-retry-and-errors.md)

## Core Rules

- Register the module with `client.Module()`.
- Build clients from `*client.Factory` rather than constructing Resty clients manually.
- Keep client behavior config-driven and request execution context-aware.
