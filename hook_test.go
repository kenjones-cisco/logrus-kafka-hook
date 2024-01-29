package logkafka

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/sirupsen/logrus"
)

var mockData = []byte("test_data")

type testReporterMock struct {
	errors []string
}

func newTestReporterMock() *testReporterMock {
	return &testReporterMock{errors: make([]string, 0)}
}

func (trm *testReporterMock) Errorf(format string, args ...interface{}) {
	trm.errors = append(trm.errors, fmt.Sprintf(format, args...))
}

func makeProducer(er mocks.ErrorReporter, vc mocks.ValueChecker) sarama.AsyncProducer {
	c := mocks.NewAsyncProducer(er, nil)
	c.ExpectInputWithCheckerFunctionAndSucceed(vc)

	return c
}

func makeErrorProducer(er mocks.ErrorReporter, vc mocks.ValueChecker, err error) sarama.AsyncProducer {
	c := mocks.NewAsyncProducer(er, nil)
	c.ExpectInputWithCheckerFunctionAndFail(vc, err)

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

func makeErrorValueChecker(err error) mocks.ValueChecker {
	return func([]byte) error {
		return err
	}
}

type simpleFmter struct{}

func (f simpleFmter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(e.Message), nil
}

func TestFire(t *testing.T) {
	trm := newTestReporterMock()

	h := New().WithFormatter(simpleFmter{}).WithProducer(makeProducer(trm, makeValueChecker(mockData))).WithTopic("success")

	entry := &logrus.Entry{
		Message: string(mockData),
		Data:    logrus.Fields{},
	}

	if err := h.Fire(entry); err != nil {
		t.Errorf("Fire() expected Fire to not return error got: %v", err)
	}

	if err := h.producer.Close(); err != nil {
		t.Error(err)
	}

	if len(trm.errors) != 0 {
		t.Errorf("Expected no errors got: %v", trm.errors)
	}
}

func TestFire_NoProducer(t *testing.T) {
	h := New()

	entry := &logrus.Entry{
		Message: string(mockData),
		Data:    logrus.Fields{},
		Time:    time.Date(0, 1, 2, 3, 4, 5, 6, time.FixedZone("", -1*60)), // cause the MarshalBinary for entry.Time to fail
	}

	if err := h.Fire(entry); err == nil {
		t.Error("Fire() expected Fire to return error got: nil")
	} else if err.Error() != "no producer defined" {
		t.Errorf("Fire() wanted: no producer defined got: %v", err)
	}
}

type failFmt struct{}

func (f failFmt) Format(e *logrus.Entry) ([]byte, error) {
	return nil, errors.New("formatting error")
}

func TestFire_FormatError(t *testing.T) {
	h := New().WithFormatter(failFmt{})

	entry := &logrus.Entry{
		Message: string(mockData),
		Data:    logrus.Fields{},
	}

	if err := h.Fire(entry); err == nil {
		t.Error("Fire() expected Fire to return error got: nil")
	} else if err.Error() != "formatting error" {
		t.Errorf("Fire() wanted: formatting error got: %v", err)
	}
}

func TestFire_PublishError(t *testing.T) {
	trm := newTestReporterMock()

	wantErr := errors.New("failed to publish message")
	h := New().WithProducer(makeErrorProducer(trm, makeErrorValueChecker(wantErr), wantErr)).WithTopic("failure")

	entry := &logrus.Entry{
		Message: string(mockData),
		Data:    logrus.Fields{},
	}

	if err := h.Fire(entry); err != nil {
		t.Errorf("Fire() expected Fire to not return error got: %v", err)
	}

	if err := h.producer.Close(); err != nil {
		t.Error(err)
	}

	if len(trm.errors) != 1 {
		t.Error("Expected to report an error")
	}
}

func TestNew(t *testing.T) {
	if h := New(); h == nil {
		t.Error("New() Expected non-nil value")
	}
}

func TestNew_With(t *testing.T) {
	h := New()
	if h.topic != "logs" {
		t.Errorf("New() want: logs got: %v", h.topic)
	}

	levels := []logrus.Level{logrus.ErrorLevel}

	oh := h.WithFormatter(simpleFmter{}).WithTopic("other").WithLevels(levels)

	if oh.topic != "other" {
		t.Errorf("New() want: other got: %v", h.topic)
	}

	if !reflect.DeepEqual(oh.Levels(), levels) {
		t.Errorf("Levels() want: %v got: %v", levels, oh.Levels())
	}

	if oh.producer != nil {
		t.Error("New() expected producer to be nil by default")
	}
}
