package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

const GoTime = "2006-01-02 15:04:05"

// GiNanaStdFormatter 自定义 formatter
type GiNanaStdFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool
}

// Format implement the Formatter interface
func (f *GiNanaStdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = GoTime
	}
	if !f.DisableTimestamp {
		data[logrus.FieldKeyTime] = entry.Time.Format(timestampFormat)
	}
	data[logrus.FieldKeyMsg] = entry.Message
	data[logrus.FieldKeyLevel] = entry.Level.String()
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if funcVal != "" {
			data[logrus.FieldKeyFunc] = funcVal
		}
		if fileVal != "" {
			data[logrus.FieldKeyFile] = fileVal
		}
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	f.toStdOut(b, fmt.Sprintf("[%s]", data["appName"]))
	f.toStdOut(b, fmt.Sprintf("[%s]", data[logrus.FieldKeyLevel]))
	if !f.DisableTimestamp {
		f.toStdOut(b, data[logrus.FieldKeyTime])
	}
	b.WriteString(" |")
	f.toStdOut(b, data[logrus.FieldKeyMsg])
	//b.WriteString(" |")
	//f.toStdOut(b, fmt.Sprintf("[%s]", data["stack"]))
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *GiNanaStdFormatter) toStdOut(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}
