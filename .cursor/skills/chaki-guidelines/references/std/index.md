# Standard Modules

Read this when the task depends on Chaki's built-in configuration and logging behavior.

## Scope

Standard modules are available by default. They are foundational for the other Chaki modules and normally do not need `app.Use(...)`.

## Read Next

- Config paths, config references, and logger usage: [config-and-logger.md](config-and-logger.md)

## Rules

- Treat config, logger, and built-in health behavior as shared infrastructure that other modules build on.
- Prefer the built-in standard-module patterns before introducing new bootstrap helpers.
