package logkafka

import (
	"crypto/tls"
	"time"

	"github.com/Shopify/sarama"
)

// SimpleProducer accepts a minimal set of configurations and creates an AsyncProducer.
func SimpleProducer(brokers []string, compression sarama.CompressionCodec, ack sarama.RequiredAcks, tls *tls.Config) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = ack
	cfg.Producer.Compression = compression
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond

	if tls != nil {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = tls
	}

	return sarama.NewAsyncProducer(brokers, cfg)
}
