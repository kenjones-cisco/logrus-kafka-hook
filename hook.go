package logkafka

import (
	"errors"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// Hook represents a logrus hook for Kafka
type Hook struct {
	formatter logrus.Formatter
	levels    []logrus.Level
	producer  sarama.AsyncProducer
	topic     string
}

// New creates a new logrus hook for Kafka
//
// Defaults:
//
//		Formatter: *logrus.TextFormatter*
//		Levels: *logrus.AllLevels*
//		Topic: *"logs"*
//
func New() *Hook {
	return &Hook{
		formatter: new(logrus.TextFormatter),
		levels:    logrus.AllLevels,
		topic:     "logs",
	}
}

// WithFormatter adds a formatter to the created Hook
func (h *Hook) WithFormatter(formatter logrus.Formatter) *Hook {
	h.formatter = formatter
	return h
}

// WithLevels adds levels to the created Hook
func (h *Hook) WithLevels(levels []logrus.Level) *Hook {
	h.levels = levels
	return h
}

// WithProducer adds a producer to the created Hook
func (h *Hook) WithProducer(producer sarama.AsyncProducer) *Hook {
	h.producer = producer

	if producer != nil {
		go func() {
			for err := range producer.Errors() {
				fmt.Fprintln(os.Stderr, "[logkafka:ERROR]", err)
			}
		}()
	}

	return h
}

// WithTopic adds a topic to the created Hook
func (h *Hook) WithTopic(topic string) *Hook {
	h.topic = topic
	return h
}

// Levels returns all log levels that are enabled for writing messages to Kafka
func (h *Hook) Levels() []logrus.Level {
	return h.levels
}

// Fire writes the entry as a message on Kafka
func (h *Hook) Fire(entry *logrus.Entry) error {
	var key sarama.Encoder

	if t, err := entry.Time.MarshalBinary(); err == nil {
		key = sarama.ByteEncoder(t)
	} else {
		key = sarama.StringEncoder(entry.Level.String())
	}

	msg, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	if h.producer == nil {
		return errors.New("no producer defined")
	}

	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: h.topic,
		Key:   key,
		Value: sarama.ByteEncoder(msg),
	}

	return nil
}
