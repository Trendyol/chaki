# Core

Read this first when the task affects app bootstrap, module composition, or more than one Chaki module.

## Read This When

- The task starts at `main.go` or overall app setup.
- The task mentions `chaki.New()`, `app.Use`, `app.Provide`, or `app.Start()`.
- The task crosses server, client, Kafka, ORM, Couchbase, or observability boundaries.

## Read Next

- App wiring and bootstrap flow: [app-composition.md](app-composition.md)
- Cross-module context rules: [context-propagation.md](context-propagation.md)

## Core Rules

- Compose applications through Chaki modules instead of hand-built service locators.
- Keep startup code thin: register modules, provide constructors, then start the app.
- Prefer config-driven module setup over hardcoded runtime values.
- If a task spans multiple modules, combine this section with the matching module references.
