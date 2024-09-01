package producer

import (
	"strings"
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/slc"
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type producerConfig struct {
	Topic   string
	Brokers string
	Writer  writerConfig
	SASL    *saslConfig
	TLS     *tlsConfig
}

type writerConfig struct {
	ReadTimeout  time.Duration
	BatchTimeout time.Duration
	BatchBytes   int64
	WriteTimeout time.Duration
	// RequiredAcks           kafka.RequiredAcks
	BatchSize              int
	WriteBackoffMax        time.Duration
	WriteBackoffMin        time.Duration
	MaxAttempts            int
	Async                  bool
	Compression            string
	AllowAutoTopicCreation bool
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

func buildLibConfig(cfg *config.Config, producerName string) (*kafkalib.ProducerConfig, error) {
	pcfg, err := config.ToStruct[producerConfig](cfg.Of(producerKey), producerName)
	if err != nil {
		return nil, err
	}

	libcfg := &kafkalib.ProducerConfig{
		Writer: kafkalib.WriterConfig{
			Topic:        pcfg.Topic,
			Brokers:      slc.Map(strings.Split(pcfg.Brokers, ","), strings.TrimSpace),
			ReadTimeout:  pcfg.Writer.ReadTimeout,
			BatchTimeout: pcfg.Writer.BatchTimeout,
			BatchBytes:   pcfg.Writer.BatchBytes,
			WriteTimeout: pcfg.Writer.WriteTimeout,
			// RequiredAcks: , TODO:
			BatchSize:       pcfg.Writer.BatchSize,
			WriteBackoffMax: pcfg.Writer.WriteBackoffMax,
			WriteBackoffMin: pcfg.Writer.WriteBackoffMin,
			MaxAttempts:     pcfg.Writer.MaxAttempts,
			Async:           pcfg.Writer.Async,
			// Compression: compress.Compression(),
			AllowAutoTopicCreation: pcfg.Writer.AllowAutoTopicCreation,
		},
	}

	if pcfg.SASL != nil {
		libcfg.SASL = &kafkalib.SASLConfig{
			Type:     kafkalib.Mechanism(pcfg.SASL.Type),
			Username: pcfg.SASL.Username,
			Password: pcfg.SASL.Password,
		}
	}

	if pcfg.TLS != nil {
		libcfg.TLS = &kafkalib.TLSConfig{
			RootCAPath:         pcfg.TLS.RootCAPath,
			IntermediateCAPath: pcfg.TLS.IntermediateCAPath,
		}
	}

	return libcfg, nil
}
