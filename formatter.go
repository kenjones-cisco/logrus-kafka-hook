package logkafka

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

// Using a pool to re-use of old entries when formatting messages.
// It is used in the Fire function.
var entryPool = sync.Pool{
	New: func() interface{} {
		return &logrus.Entry{}
	},
}

// copyEntry copies the entry `e` to a new entry and then adds all the fields in `fields` that are missing in the new entry data.
// It uses `entryPool` to re-use allocated entries.
func copyEntry(e *logrus.Entry, fields logrus.Fields) *logrus.Entry {
	ne := entryPool.Get().(*logrus.Entry)
	ne.Message = e.Message
	ne.Level = e.Level
	ne.Time = e.Time
	ne.Data = logrus.Fields{}

	for k, v := range fields {
		ne.Data[k] = v
	}

	for k, v := range e.Data {
		ne.Data[k] = v
	}

	return ne
}

// releaseEntry puts the given entry back to `entryPool`. It must be called if copyEntry is called.
func releaseEntry(e *logrus.Entry) {
	entryPool.Put(e)
}

// StructuredFormatter represents a structured format.
// It has logrus.Formatter which formats the entry and logrus.Fields which
// are added to the JSON message if not given in the entry data.
//
//		**Note:** use the `DefaultFormatter` function to create a default StructuredFormatter.
//
type StructuredFormatter struct {
	logrus.Formatter
	logrus.Fields
}

var (
	logFields   = logrus.Fields{"@version": "1", "type": "log"}
	logFieldMap = logrus.FieldMap{
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyMsg:  "message",
	}
)

// DefaultFormatter returns a default structured formatter:
// A JSON format with "@version" set to "1" (unless set differently in `fields`,
// "type" to "log" (unless set differently in `fields`),
// "@timestamp" to the log time and "message" to the log message.
//
//		**Note:** to set a different configuration use the `StructuredFormatter` structure.
//
func DefaultFormatter(fields logrus.Fields) logrus.Formatter {
	for k, v := range logFields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}

	return StructuredFormatter{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: RFC3339NanoFixed,
			FieldMap:        logFieldMap,
		},
		Fields: fields,
	}
}

// Format formats an entry to a structured format according to the given Formatter and Fields.
//
//		**Note:** the given entry is copied and not changed during the formatting process.
//
func (f StructuredFormatter) Format(e *logrus.Entry) ([]byte, error) {
	ne := copyEntry(e, f.Fields)

	dataBytes, err := f.Formatter.Format(ne)
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	releaseEntry(ne)

	return dataBytes, err
}
