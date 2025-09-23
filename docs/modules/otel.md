# OTEL Module

The OTEL module provides OpenTelemetry integration for tracing and metrics in Chaki applications.

## Submodules

- **Client**: Tracing for HTTP clients.
- **Server**: Tracing for HTTP servers.
- **Kafka**: Tracing for Kafka producers and consumers.

## Usage

```go
app.Use(otel.Module(
  otelclient.WithClient(),
  otelserver.WithServer(),
  otelkafka.WithKafka(),
))
```

Have to provide an otel init function that initiates your otel exporting mechanism

See option.go for more details.
