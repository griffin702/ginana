package log

import "io"

var (
	log   Logger
)

func GetOutFile() (out io.Writer) {
	return log.GetOutFile()
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func PrintErrf(format string, args ...interface{}) {
	log.PrintErrf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
