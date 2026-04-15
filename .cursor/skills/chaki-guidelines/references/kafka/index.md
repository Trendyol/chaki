# Kafka Module

Use this page when wiring Kafka producers, consumers, or Kafka-specific config.

## Read Next

- Producer and consumer usage patterns: [producers-and-consumers.md](producers-and-consumers.md)

## Core Rules

- Keep producer and consumer wiring config-driven.
- Use caller context for producers and consumed message context for downstream handling.
- Keep reused producer names, consumer names, and topic-related identifiers in constants.
- Chaki also supports batch consumers; if a task is explicitly batch-oriented, validate the implementation against the batch consumer APIs instead of assuming the single-message consumer shape.

## Config Quick Map

```yaml
kafka:
  producer:
    order-events:
      brokers: "broker1:9092,broker2:9092"

  consumer:
    order-created-consumer:
      brokers: "broker1:9092,broker2:9092"
      topic: "order.created.v1"
      groupId: "order.listener"
```

- `kafka.producer.<name>` must match `producer.Factory.Get("<name>")`.
- `kafka.consumer.<name>` must match `consumer.New("<name>")` or `consumer.NewFn("<name>", fn)`.
- Producer `brokers` is forgiving about spaces after commas, but consumer `brokers` is not.
- Consumer `brokers` is a comma-separated string, and it is safest to keep it compact without spaces after commas.
- `retryConfiguration.brokers` should be documented as a YAML list.
