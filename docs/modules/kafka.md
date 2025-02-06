# Kafka Module

The Kafka module is a robust and flexible solution built on top of [kafka-konsumer](https://github.com/Trendyol/kafka-konsumer) and [kafka-go](https://github.com/segmentio/kafka-go). It simplifies the process of interacting with Apache Kafka by enabling you to create configuration-driven consumers and producers with minimal effort, saving both time and complexity in setup and usage. Whether you're working with simple or complex Kafka configurations, this module offers a streamlined approach to managing your Kafka-based infrastructure.


## Configuration

The Kafka module supports nearly all configuration options provided by kafka-konsumer. You can define these configuration in consumer and producer configuration sections. Below is an example of how to define consumer and producer configurations in YAML format:


```yaml

kafka:
  consumer:
    <CONSUMER_NAME_1>:
        topic: string          # The topic to consume messages from
        brokers: string        # Comma-separated list of broker addresses
        groupId: string        # Consumer group ID for coordinating message consumption
        clientId: string       # Client ID used for identifying the consumer
        # Additional consumer configurations can be specified here.
        # Refer to the library documentation for the full list of options:
        # https://github.com/Trendyol/kafka-konsumer?tab=readme-ov-file#configurations

  producer:
    <PRODUCER_NAME_1>:
        topic: string          # Optional: Default topic for producing messages
        brokers: string        # Comma-separated list of broker addresses
        # Additional producer configurations can be specified here.
        # Refer to the library documentation for the full list of options:
        # https://github.com/Trendyol/kafka-konsumer?tab=readme-ov-file#configurations 


```


## Example

Here is an example configuration and usage for both a consumer and a producer, including settings for error retry, SASL authentication, and TLS configurations.

#### config.yaml
```yaml

BROKERS: "kafka-brokers:9090, kafka-brokers:9091"

kafka:
    consumer:
        foo:
            topic: "example-topic"
            brokers: ${this:BROKERS}
            groupId: "foo-consumer-group"
            clientId: "foo-app-client"
            sasl.type: scram
            sasl.username: "username"
            sasl.password: "password" # you can read from another or reference. ${env:sasl.password}, ${secret:sasl.password}
            tls.rootCaPath: ./cert-path/cert.pem
            tls.intermediateCaPath:  ./cert-path/cert.pem
            retryEnabled: true
            retryConfiguration:
                topic: "example-topic-retry"
                startTimeCron: "* * * * *"
                deadLetterTopic: "example-topic-err"    

    producer:
        foobar:
            topic: "example-topic"
            brokers: ${this:BROKERS}
```

#### service.go
```go
import (
    "github.com/Trendyol/chaki/modules/kafka/producer"
)

type Service struct {
    p producer.Producer
}

func NewService(pf producer.Factory) *Service {
    return &Service{
        p: pf.Get("foobar")
    }
}

func (s *Service) Notify(ctx context.Context) error {
    return s.p.Produce(ctx, producer.Message{
        Value: []byte("hello world!")
    })
}

```

#### consumer.go

```go
import (
    "github.com/Trendyol/chaki/modules/kafka/consumer"
    "github.com/Trendyol/chaki/logger"
)

type Consumer struct {
    consumer.Consumer // composition
}

func NewConsumer() *Consumer {
    return &Service{
        Consumer: consumer.New("foo"),
    }
}

func (c *Consumer) Consume(ctx context.Context, msg *consumer.Message) error {
    logger.From(msg.Context).Info("message consumed")
    // Handle message
    return nil
}

```







