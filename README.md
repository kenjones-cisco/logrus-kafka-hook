# logkafka

## Status

[![GoDoc][1]][2]
[![GoCard][3]][4]
[![coverage][5]][6]
[![Build Status][7]][8]

[1]: https://godoc.org/github.com/kenjones-cisco/logrus-kafka-hook?status.svg
[2]: https://godoc.org/github.com/kenjones-cisco/logrus-kafka-hook
[3]: https://goreportcard.com/badge/kenjones-cisco/logrus-kafka-hook
[4]: https://goreportcard.com/report/github.com/kenjones-cisco/logrus-kafka-hook
[5]: http://gocover.io/_badge/github.com/kenjones-cisco/logrus-kafka-hook
[6]: http://gocover.io/github.com/kenjones-cisco/logrus-kafka-hook
[7]: https://travis-ci.org/kenjones-cisco/logrus-kafka-hook.svg?branch=master
[8]: https://travis-ci.org/kenjones-cisco/logrus-kafka-hook


`import "github.com/kenjones-cisco/logrus-kafka-hook"`

Package logkafka provides a Logrus Hook implementation for publishing log messages on a Kafka Topic.

logkafka uses a builder pattern to configure the Hook. This provides simplicity and flexibility when configuring a Hook.
There are convenience functions for creating a sarama.AsyncProducer and logrus.Formatter.

### Basic Usage

```go
// use SimpleProducer to create an AsyncProducer
producer := logkafka.SimpleProducer([]string{"127.0.0.1:9092"}, sarama.CompressionSnappy, sarama.WaitForLocal, nil)

// use DefaultFormatter to create a JSONFormatter with pre-defined fields or override with any fields
// create the Hook and use the builder functions to apply configurations
hook := logkafka.New().WithFormatter(logkafka.DefaultFormatter(logrus.Fields{"type": "myappName"})).WithProducer(producer)

// create a new logger and add the hook
log := logrus.New()
log.Hooks.Add(hook)

log.Debug("Hello World")
```

Example of formatted message published to Kafka

```json
{


	"@timestamp": "2018-04-20T04:03:00Z",
	"@version": "1",
	"level": "info",
	"message": "Hello World",
	"type": "myappName"

}
```
