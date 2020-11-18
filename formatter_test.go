package logkafka

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestDefaultFormatterWithFields(t *testing.T) {
	format := DefaultFormatter(logrus.Fields{"ID": 123})

	entry := &logrus.Entry{
		Message: "msg1",
		Data:    logrus.Fields{"f1": "bla"},
	}

	res, err := format.Format(entry)
	if err != nil {
		t.Errorf("expected format to not return error: %s", err)
	}

	expected := []string{
		"f1\":\"bla\"",
		"ID\":123",
		"message\":\"msg1\"",
	}

	for _, exp := range expected {
		if !strings.Contains(string(res), exp) {
			t.Errorf("expected to have '%s' in '%s'", exp, string(res))
		}
	}
}

func TestDefaultFormatterWithEmptyFields(t *testing.T) {
	now := time.Now()
	formatter := DefaultFormatter(logrus.Fields{})

	entry := &logrus.Entry{
		Message: "message bla bla",
		Level:   logrus.DebugLevel,
		Time:    now,
		Data: logrus.Fields{
			"Key1": "Value1",
		},
	}

	res, err := formatter.Format(entry)
	if err != nil {
		t.Errorf("expected Format not to return error: %s", err)
	}

	expected := []string{
		"\"message\":\"message bla bla\"",
		"\"level\":\"debug\"",
		"\"Key1\":\"Value1\"",
		"\"@version\":\"1\"",
		"\"type\":\"log\"",
		fmt.Sprintf("\"@timestamp\":\"%s\"", now.Format(RFC3339NanoFixed)),
	}

	for _, exp := range expected {
		if !strings.Contains(string(res), exp) {
			t.Errorf("expected to have '%s' in '%s'", exp, string(res))
		}
	}
}

func TestLogFieldsNotOverridden(t *testing.T) {
	_ = DefaultFormatter(logrus.Fields{"user1": "11"})

	if _, ok := logFields["user1"]; ok {
		t.Errorf("expected user1 to not be in logFields: %#v", logFields)
	}
}
