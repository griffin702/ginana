package log

import (
	"github.com/griffin702/ginana/library/log/hook"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/ulricqin/goutils/filetool"
	"io"
	"os"
	"time"
)

var (
	_level        logrus.Level
	_path         string
	_maxAge       uint
	_rotationTime uint
)

type logger struct {
	log     *logrus.Logger
	logType int
}

func (l *logger) GetOutFile() (out io.Writer) {
	return l.log.Out
}

func (l *logger) isStdOut() bool {
	return l.logType == 0
}

func (l *logger) NewStdOut() {
	cli := logrus.New()
	cli.SetLevel(_level)
	cli.Out = os.Stdout
	cli.Formatter = &GiNanaStdFormatter{}
	cli.AddHook(&hook.DefaultFieldHook{})
	cli.AddHook(&hook.LineHook{})
	l.log = cli
	return
}

func (l *logger) NewFile() (cf func()) {
	cli := logrus.New()
	cli.SetLevel(_level)
	out, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	cf = func() {
		if out != nil {
			_ = out.Close()
		}
	}
	_ = filetool.InsureDir(_path)
	logWriter, err := rotatelogs.New(
		_path+"/log-%Y-%m-%d-%H-%M.log",
		//rotatelogs.WithLinkName(d.path),
		rotatelogs.WithMaxAge(time.Duration(_maxAge)*time.Hour),             // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(_rotationTime)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		return
	}
	cli.Out = logWriter
	cli.Formatter = &GiNanaStdFormatter{
		DisableColors: true,
	}
	//writeMap := lfshook.WriterMap{
	//	logrus.DebugLevel: logWriter,
	//	logrus.InfoLevel:  logWriter,
	//	logrus.WarnLevel:  logWriter,
	//	logrus.ErrorLevel: logWriter,
	//	logrus.FatalLevel: logWriter,
	//	logrus.PanicLevel: logWriter,
	//}
	//lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{DisableColors: true})
	cli.AddHook(&hook.DefaultFieldHook{})
	cli.AddHook(&hook.LineHook{})
	//cli.AddHook(lfHook)
	l.log = cli
	return
}

func (l *logger) Print(args ...interface{}) {
	l.log.Print(args...)
}

func (l *logger) Println(args ...interface{}) {
	l.log.Println(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}
