package log

import (
	"io"
)

type Logger interface {
	GetOutFile() (out io.Writer)
	isStdOut() bool
	Print(args ...interface{})
	Println(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

func init() {
	_level = 5
	_path = "../logs"
	_maxAge = 7 * 24
	_rotationTime = 24
}

func Init() func() {
	stdOut := new(logger)
	stdOut.logType = 0
	stdOut.NewStdOut()
	logRus := new(logger)
	logRus.logType = 1
	cf := logRus.NewFile()
	var hs Loggers
	hs = append(hs, stdOut, logRus)
	log = &hs
	return cf
}

type Loggers []Logger

func (h Loggers) GetOutFile() (out io.Writer) {
	var wr []io.Writer
	for _, l := range h {
		if o := l.GetOutFile(); o != nil {
			wr = append(wr, o)
		}
	}
	if len(wr) > 0 {
		out = io.MultiWriter(wr...)
		return
	}
	return nil
}

func (h Loggers) isStdOut() bool {
	return false
}

func (h Loggers) Print(args ...interface{}) {
	for _, l := range h {
		if l.isStdOut() {
			l.Print(args...)
		}
	}
}

func (h Loggers) Println(args ...interface{}) {
	for _, l := range h {
		if l.isStdOut() {
			l.Println(args...)
		}
	}
}

func (h Loggers) Error(args ...interface{}) {
	for _, l := range h {
		l.Error(args...)
	}
}

func (h Loggers) Warn(args ...interface{}) {
	for _, l := range h {
		l.Warn(args...)
	}
}

func (h Loggers) Info(args ...interface{}) {
	for _, l := range h {
		l.Info(args...)
	}
}

func (h Loggers) Debug(args ...interface{}) {
	for _, l := range h {
		l.Debug(args...)
	}
}

func (h Loggers) Printf(format string, args ...interface{}) {
	for _, l := range h {
		if l.isStdOut() {
			l.Printf(format, args...)
		}
	}
}

func (h Loggers) Errorf(format string, args ...interface{}) {
	for _, l := range h {
		l.Errorf(format, args...)
	}
}

func (h Loggers) Warnf(format string, args ...interface{}) {
	for _, l := range h {
		l.Warnf(format, args...)
	}
}

func (h Loggers) Infof(format string, args ...interface{}) {
	for _, l := range h {
		l.Infof(format, args...)
	}
}

func (h Loggers) Debugf(format string, args ...interface{}) {
	for _, l := range h {
		l.Debugf(format, args...)
	}
}
