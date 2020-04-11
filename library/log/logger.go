package log

import "io"

var (
	log Logger
)

func GetLogger() Logger {
	return log
}

func GetOutFile() (out io.Writer) {
	return log.GetOutFile()
}

func Print(args ...interface{}) {
	log.Print(args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.Infof(format, args...)
}
