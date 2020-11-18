package logkafka

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/Shopify/sarama"
)

const flushFrequency = 500 * time.Millisecond

// ErrNoProducer is used when hook Fire method is called and no producer was configured.
var ErrNoProducer = errors.New("no producer defined")

// SimpleProducer accepts a minimal set of configurations and creates an AsyncProducer.
func SimpleProducer(brokers []string, compression sarama.CompressionCodec, ack sarama.RequiredAcks, tlscfg *tls.Config) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = ack
	cfg.Producer.Compression = compression
	cfg.Producer.Flush.Frequency = flushFrequency

	if tlscfg != nil {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = tlscfg
	}

	return sarama.NewAsyncProducer(brokers, cfg)
}
