# Producers And Consumers

Use this page when implementing Kafka business flows.

## Config Rules

- Configure producers under `kafka.producer.<name>`.
- Configure consumers under `kafka.consumer.<name>`.
- SASL, TLS, retry, dead-letter, and retry-topic behavior belong in config, not in business handlers.

## Producer Config Example

```yaml
kafka:
  producer:
    order-events:
      brokers: "broker1:9092,broker2:9092"
      tls:
        rootCAPath: "/path/to/root-ca.pem"
        intermediateCAPath: "/path/to/intermediate-ca.pem"
      sasl:
        username: ${kafka-auth:username}
        password: ${kafka-auth:password}
        type: scram
```

- The producer name `order-events` must match `producer.Factory.Get("order-events")`.
- `brokers` is a comma-separated string in producer config.
- `topic` is optional in producer config. If it is not configured, set `producer.Message.Topic` when publishing.

## Producer Pattern

- Inject `producer.Factory`.
- Get the producer with `factory.Get(producerName)`.
- Produce messages with the caller context so tracing and metadata flow with the event.

## Consumer Pattern

- Compose `consumer.Consumer` into the module-specific consumer type.
- Construct it with `consumer.New(consumerName)` when the type will provide its own `Consume` implementation.
- For simple handlers, `consumer.NewFn(consumerName, fn)` is the lighter-weight option.
- Implement `Consume(msg *consumer.Message) error`.
- Prefer `msg.Context` when logging or calling downstream dependencies from consumed messages.

## Consumer Config Example

```yaml
kafka:
  consumer:
    order-created-consumer:
      brokers: "broker1:9092,broker2:9092"
      topic: "order.created.v1"
      groupId: "order.listener"
      tls:
        rootCAPath: "/path/to/root-ca.pem"
        intermediateCAPath: "/path/to/intermediate-ca.pem"
      sasl:
        username: ${kafka-auth:username}
        password: ${kafka-auth:password}
        type: scram
```

- The consumer name `order-created-consumer` must match `consumer.New("order-created-consumer")`.
- `groupId` belongs on the consumer config, not in handler code.
- Use reference files for credential placeholders instead of embedding usernames or passwords in YAML.
- Keep the main consumer `brokers` string compact, such as `"broker1:9092,broker2:9092"`, because the framework does not trim spaces after splitting it.

## Retry Example

```yaml
kafka:
  consumer:
    retrying-consumer:
      brokers: "broker1:9092,broker2:9092"
      topic: "seller.questionanswer.created.v1"
      groupId: "seller.questionanswer.listener"
      retryEnabled: true
      retryConfiguration:
        startTimeCron: "*/1 * * * *"
        topic: "seller.questionanswer.created.v1.retry"
        deadLetterTopic: "seller.questionanswer.created.v1.dlq"
        brokers:
          - "broker1:9092"
          - "broker2:9092"
        maxRetry: 3
        workDuration: 60
```

- Keep retry topics and dead-letter topics explicit in config.
- Prefer a YAML list for `retryConfiguration.brokers` because the framework struct expects `[]string`.

## Batch Example

```yaml
kafka:
  consumer:
    batch-consumer:
      brokers: "broker1:9092,broker2:9092"
      topic: "inventory.changed.v1"
      groupId: "inventory.listener"
      batchConfiguration:
        messageGroupLimit: 100
        messageGroupByteSizeLimit: 1048576
```

- Batch settings live under the same consumer entry that owns the topic and group.
- If `batchConfiguration` is enabled, use the batch-consumer path such as `consumer.NewBatchConsumer(...)` or `consumer.NewBatchConsumerFn(...)` instead of the single-message path.

## Design Guidance

- Keep serialization, topic selection, and business handling explicit and readable.
- When consumer or producer names are reused across config and code, define constants instead of repeating literals.
