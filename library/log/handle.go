package log

import (
	"io"
)

type Logger interface {
	GetOutFile() (out io.Writer)
	isStdOut() bool
	Printf(format string, args ...interface{})
	PrintErrf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

func init() {
	_level = 5
	_path = "../logs/log"
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

func (h Loggers) Printf(format string, args ...interface{}) {
	for _, l := range h {
		if l.isStdOut() {
			l.Printf(format, args...)
		}
	}
}

func (h Loggers) PrintErrf(format string, args ...interface{}) {
	for _, l := range h {
		l.Printf(format, args...)
	}
}

func (h Loggers) Info(args ...interface{}) {
	for _, l := range h {
		l.Info(args...)
	}
}

func (h Loggers) Infof(format string, args ...interface{}) {
	for _, l := range h {
		l.Infof(format, args...)
	}
}

func (h Loggers) Error(args ...interface{}) {
	for _, l := range h {
		l.Error(args...)
	}
}

func (h Loggers) Errorf(format string, args ...interface{}) {
	for _, l := range h {
		l.Errorf(format, args...)
	}
}
