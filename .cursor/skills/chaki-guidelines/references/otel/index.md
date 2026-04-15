# OTEL Module

Use this page when the task wires OpenTelemetry into Chaki modules.

## Read Next

- Module and submodule wiring: [instrumentation-wiring.md](instrumentation-wiring.md)

## Core Rules

- Enable only the OTEL submodules that correspond to the Chaki modules the app actually uses.
- Keep exporter initialization separate from business logic.
- Combine this page with the server, client, or Kafka references when instrumentation changes touch those modules too.
