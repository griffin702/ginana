package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

const GoTime = "2006-01-02 15:04:05"

var (
	//greenBg      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	//whiteBg      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	//yellowBg     = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	//redBg        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	//blueBg       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	//magentaBg    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	//cyanBg       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	green   = string([]byte{27, 91, 51, 50, 109})
	white   = string([]byte{27, 91, 51, 55, 109})
	yellow  = string([]byte{27, 91, 51, 51, 109})
	red     = string([]byte{27, 91, 51, 49, 109})
	blue    = string([]byte{27, 91, 51, 52, 109})
	magenta = string([]byte{27, 91, 51, 53, 109})
	cyan    = string([]byte{27, 91, 51, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// GiNanaStdFormatter 自定义 formatter
type GiNanaStdFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool
	DisableColors    bool
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
	var nameColor, levelColor, resetColor string
	levelValue := data[logrus.FieldKeyLevel]
	if !f.DisableColors {
		nameColor = f.NameColor()
		levelColor = f.LevelColor(levelValue)
		resetColor = f.ResetColor()
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	f.toStdOut(b, fmt.Sprintf("%s[%s]%s", nameColor, data["appName"], resetColor))
	f.toStdOut(b, fmt.Sprintf("%s[%s]%s", levelColor, levelValue, resetColor))
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

func (f *GiNanaStdFormatter) ResetColor() string {
	return reset
}

func (f *GiNanaStdFormatter) NameColor() string {
	return magenta
}

func (f *GiNanaStdFormatter) LevelColor(level interface{}) string {
	switch level {
	case "disable":
		return white
	case "fatal":
		return cyan
	case "error":
		return red
	case "warn":
		return yellow
	case "info":
		return green
	case "debug":
		return blue
	default:
		return reset
	}
}
