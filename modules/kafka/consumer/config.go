package consumer

import (
	"strings"
	"time"

	"github.com/Trendyol/chaki/config"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
	segmentio "github.com/segmentio/kafka-go"
)

type consumerConfig struct {
	Brokers  string
	GroupID  string
	Topic    string
	ClientID string
	Reader   readerConfig
	SASL     *saslConfig
	TLS      *tlsConfig
}

type readerConfig struct {
	Partition              int
	QueueCapacity          int
	MinBytes               int
	MaxBytes               int
	MaxWait                time.Duration
	ReadBatchTimeout       time.Duration
	ReadLagInterval        time.Duration
	GroupBalancers         string
	HeartbeatInterval      time.Duration
	CommitInterval         time.Duration
	PartitionWatchInterval time.Duration
	WatchPartitionChanges  bool
	SessionTimeout         time.Duration
	RebalanceTimeout       time.Duration
	JoinGroupBackoff       time.Duration
	RetentionTime          time.Duration
	StartOffset            int64
	ReadBackoffMin         time.Duration
	ReadBackoffMax         time.Duration
	IsolationLevel         segmentio.IsolationLevel
	MaxAttempts            int
	OffsetOutOfRangeError  bool
}

type tlsConfig struct {
	RootCAPath         string
	IntermediateCAPath string
}

type saslConfig struct {
	Type     string
	Username string
	Password string
}

func buildLibConfig(cfg *config.Config, consumer Consumer) (*kafkalib.ConsumerConfig, error) {
	kcfg, err := config.ToStruct[consumerConfig](cfg.Of("kafka.consumer"), consumer.Name())
	if err != nil {
		return nil, err
	}

	libcfg := &kafkalib.ConsumerConfig{
		ClientID:  kcfg.ClientID,
		ConsumeFn: consumer.Consume,
		Reader: kafkalib.ReaderConfig{
			Brokers:                splitByComma(kcfg.Brokers),
			Topic:                  kcfg.Topic,
			GroupID:                kcfg.GroupID,
			Partition:              kcfg.Reader.Partition,
			QueueCapacity:          kcfg.Reader.QueueCapacity,
			MinBytes:               kcfg.Reader.MinBytes,
			MaxBytes:               kcfg.Reader.MaxBytes,
			MaxWait:                kcfg.Reader.MaxWait,
			ReadBatchTimeout:       kcfg.Reader.ReadBatchTimeout,
			ReadLagInterval:        kcfg.Reader.ReadLagInterval,
			GroupBalancers:         nil, // TODO:
			HeartbeatInterval:      kcfg.Reader.HeartbeatInterval,
			CommitInterval:         kcfg.Reader.CommitInterval,
			PartitionWatchInterval: kcfg.Reader.PartitionWatchInterval,
			WatchPartitionChanges:  kcfg.Reader.WatchPartitionChanges,
			SessionTimeout:         kcfg.Reader.SessionTimeout,
			RebalanceTimeout:       kcfg.Reader.ReadBatchTimeout,
			JoinGroupBackoff:       kcfg.Reader.JoinGroupBackoff,
			RetentionTime:          kcfg.Reader.RetentionTime,
			StartOffset:            kcfg.Reader.StartOffset,
			ReadBackoffMin:         kcfg.Reader.ReadBackoffMin,
			ReadBackoffMax:         kcfg.Reader.ReadBackoffMax,
			IsolationLevel:         kcfg.Reader.IsolationLevel,
			MaxAttempts:            kcfg.Reader.MaxAttempts,
			OffsetOutOfRangeError:  kcfg.Reader.OffsetOutOfRangeError,
		},
	}

	if kcfg.SASL != nil {
		libcfg.SASL = &kafkalib.SASLConfig{
			Type:     getSaslMechanism(kcfg.SASL.Type),
			Password: kcfg.SASL.Password,
			Username: kcfg.SASL.Username,
		}
	}

	if kcfg.TLS != nil {
		libcfg.TLS = &kafkalib.TLSConfig{
			RootCAPath:         kcfg.TLS.RootCAPath,
			IntermediateCAPath: kcfg.TLS.IntermediateCAPath,
		}
	}

	return libcfg, nil
}

func splitByComma(s string) []string {
	return strings.Split(s, ",")
}

func getSaslMechanism(s string) kafkalib.Mechanism {
	if s == "" {
		return kafkalib.MechanismScram
	}
	return kafkalib.Mechanism(s)
}
