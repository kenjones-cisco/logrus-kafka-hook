package logkafka_test

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"

	logkafka "github.com/kenjones-cisco/logrus-kafka-hook"
)

func ExampleNew() {
	// use SimpleProducer to create an AsyncProducer
	producer, err := logkafka.SimpleProducer([]string{"127.0.0.1:9092"}, sarama.CompressionSnappy, sarama.WaitForLocal, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	// use DefaultFormatter to create a JSONFormatter with pre-defined fields or override with any fields
	// create the Hook and use the builder functions to apply configurations
	hook := logkafka.New().WithFormatter(
		logkafka.DefaultFormatter(logrus.Fields{"type": "myappName"})).WithProducer(producer)

	// create a new logger and add the hook
	log := logrus.New()
	log.Hooks.Add(hook)

	log.Debug("Hello World")
}
