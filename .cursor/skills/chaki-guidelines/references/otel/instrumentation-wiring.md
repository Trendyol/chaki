# Instrumentation Wiring

Use this page when enabling or modifying OpenTelemetry support.

## Wiring Pattern

- Register the module with `otel.Module(...)`.
- Add only the needed submodule options from the OTEL subpackages, such as `otelclient.WithClient()`, `otelserver.WithServer()`, and `otelkafka.WithKafka()`, then pass them into `otel.Module(...)`.
- Provide the OTEL initialization and exporter setup required by the target environment.

## Design Guidance

- Match instrumentation scope to actual module usage instead of enabling every option by default.
- Preserve request context across server, client, Kafka, and persistence boundaries so spans can join correctly.
- If a task changes both instrumentation and transport or data modules, read those module references too.
