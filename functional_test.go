package logkafka_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/sirupsen/logrus"

	logkafka "github.com/kenjones-cisco/logrus-kafka-hook"
)

var (
	mockData   = []byte("test_data")
	warnLevels = []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
)

type testReporterMock struct {
	errors []string
}

func newTestReporterMock() *testReporterMock {
	return &testReporterMock{errors: make([]string, 0)}
}

func (trm *testReporterMock) Errorf(format string, args ...interface{}) {
	trm.errors = append(trm.errors, fmt.Sprintf(format, args...))
}

func makeProducer(er mocks.ErrorReporter, vc ...mocks.ValueChecker) sarama.AsyncProducer {
	c := mocks.NewAsyncProducer(er, nil)

	for _, check := range vc {
		if check != nil {
			c.ExpectInputWithCheckerFunctionAndSucceed(check)
		} else {
			c.ExpectInputAndSucceed()
		}
	}

	return c
}

func makeValueChecker(val []byte) mocks.ValueChecker {
	return func(v []byte) error {
		if !bytes.Equal(val, v) {
			return fmt.Errorf("Expected: %s, got: %s", string(val), string(v))
		}

		return nil
	}
}

type simpleFmter struct{}

func (f simpleFmter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(e.Message), nil
}

func Test_EntryPublish(t *testing.T) {
	bufferOut := bytes.NewBufferString("")

	log := logrus.New()
	log.Out = bufferOut

	trm := newTestReporterMock()
	pr := makeProducer(trm, nil)
	h := logkafka.New().WithFormatter(logkafka.DefaultFormatter(logrus.Fields{"NICKNAME": ""})).WithProducer(pr)

	log.Hooks.Add(h)

	log.Info(string(mockData))

	if strings.Contains(bufferOut.String(), "NICKNAME\":") {
		t.Errorf("expected main logrus message to not have '%s': %#v", "NICKNAME\":", bufferOut.String())
	}

	if err := pr.Close(); err != nil {
		t.Error(err)
	}

	if len(trm.errors) != 0 {
		t.Errorf("Expected no errors got: %v", trm.errors)
	}
}

func Test_EntryWithLevels(t *testing.T) {
	log := logrus.New()

	trm := newTestReporterMock()
	pr := makeProducer(trm, makeValueChecker(mockData))
	h := logkafka.New().WithLevels(warnLevels).WithFormatter(simpleFmter{}).WithProducer(pr)
	log.Hooks.Add(h)

	log.Debug(string(mockData))
	log.Warn(string(mockData))

	if err := pr.Close(); err != nil {
		t.Error(err)
	}

	if len(trm.errors) != 0 {
		t.Errorf("Expected no errors got: %v", trm.errors)
	}
}
