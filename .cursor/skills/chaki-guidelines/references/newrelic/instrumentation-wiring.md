# Instrumentation Wiring

Use this page when wiring New Relic into Chaki applications.

## Wiring Pattern

- Register the module with `newrelic.Module(...)`.
- Add only the submodule options the application needs, such as server or ORM integration.
- Keep New Relic SDK-specific tuning inside module options such as `WithNewrelicConfigOptions(...)`.

## Config Guidance

- Configure the module under the `newrelic` config branch, for example `newrelic.enabled`, `newrelic.logenabled`, `newrelic.appname`, and `newrelic.license`.
- Keep New Relic configuration externalized and scoped to application bootstrap.

## Design Guidance

- Keep monitoring setup in application composition rather than scattering it across services.
- Pair New Relic module wiring with the corresponding server, client, Kafka, or ORM references when the task spans both sides.
