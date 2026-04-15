---
name: chaki-guidelines
description: Chaki framework conventions for app composition, config-first module wiring, server controllers and routes, HTTP clients, Kafka, ORM, Couchbase, observability, Swagger, context propagation, and middleware. Use when working in this repository or when the user mentions Chaki modules, app.Use, app.Provide, controllers, routes, validation, clients, Kafka, ORM, Couchbase, OTEL, New Relic, Swagger, context, or middlewares.
---

# Chaki Guidelines

Use this skill for Chaki-specific work. Keep the main skill lightweight and load only the references that match the task.

## When To Activate

- Use this skill when changing code in this repository.
- Use this skill when the user asks how to wire or use a Chaki module.
- Use this skill when the task mentions `app.Use`, `app.Provide`, controllers, routes, validation, HTTP clients, Kafka, ORM, Couchbase, OTEL, New Relic, Swagger, context propagation, or middleware composition.

## How It Works

1. Identify the active module or cross-cutting concern.
2. Read `references/core/index.md` if the task touches app setup, multiple modules, or overall framework behavior.
3. Read the relevant module `index.md`.
4. From that module index, read only the leaf reference files needed for the task.
5. Follow existing Chaki patterns instead of introducing a new structure.

## Reference Routing

Start with these entry points:

| Task area | Read first |
| --- | --- |
| App bootstrapping, `chaki.New()`, `app.Use`, `app.Provide` | `references/core/index.md` |
| Standard config and logging behavior | `references/std/index.md` |
| Fiber server, controllers, routes, validation, middleware, error handling | `references/server/index.md` |
| REST clients, request wrappers, retries, circuit breaker, error decoding | `references/client/index.md` |
| Kafka producers and consumers | `references/kafka/index.md` |
| ORM, transactions, repositories | `references/orm/index.md` |
| Couchbase cluster usage and tracing | `references/couchbase/index.md` |
| OpenTelemetry module wiring | `references/otel/index.md` |
| New Relic module wiring | `references/newrelic/index.md` |
| Swagger generation and access restrictions | `references/swagger/index.md` |

## Working Rules

- Read only the references needed for the current task.
- Prefer Chaki's config-first, module-based composition over ad hoc wiring.
- Preserve `context.Context` across module boundaries.
- When names such as client keys, producer names, consumer names, or controller prefixes are reused, extract them into constants instead of scattering string literals.
- Keep repository changes aligned with the patterns already used by Chaki docs and examples.

## Output Expectations

- Use Chaki module APIs and naming consistently.
- Keep guidance and code changes scoped to the relevant module.
- If a task spans multiple modules, combine the corresponding module references rather than relying on one page alone.
