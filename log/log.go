package log

import (
	"context"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zmon-deploy/zmon-common-go/contextlog"
	"os"
	"time"
)

type LoggerFactory interface {
	GetLogger(component string) Logger
}

type loggerFactory struct {
}

func NewLoggerFactory() (LoggerFactory, error) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&nested.Formatter{
		NoColors: true,
	})

	return &loggerFactory{}, nil
}

func (f *loggerFactory) GetLogger(component string) Logger {
	return newLogger(component)
}

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Stopwatch(string) func()
	Stopwatchf(string, ...interface{}) func()
	WithCtx(ctx context.Context) Logger
}

type logger struct {
	component string
	ctx       context.Context
}

func newLogger(component string) Logger {
	return &logger{
		component: component,
		ctx:       nil,
	}
}

func withCtx(component string, ctx context.Context) Logger {
	return &logger{
		component: component,
		ctx:       ctx,
	}
}

func (l *logger) generateMessage(args ...interface{}) []interface{} {
	var arr []interface{}

	arr = append(arr, fmt.Sprintf("[%s] ", l.component))
	arr = append(arr, args...)

	if l.ctx != nil {
		fields := contextlog.FieldsFromContextAsLine(l.ctx)
		if len(fields) > 0 {
			arr = append(arr, ", contexts: ", fields)
		}
	}

	return arr
}

func (l *logger) Info(args ...interface{}) {
	logrus.Info(l.generateMessage(args...)...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	logrus.Info(l.generateMessage(fmt.Sprintf(format, args...))...)
}

func (l *logger) Debug(args ...interface{}) {
	logrus.Debug(l.generateMessage(args...)...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	logrus.Debug(l.generateMessage(fmt.Sprintf(format, args...))...)
}

func (l *logger) Warn(args ...interface{}) {
	logrus.Warn(l.generateMessage(args...)...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	logrus.Warn(l.generateMessage(fmt.Sprintf(format, args...))...)
}

func (l *logger) Error(args ...interface{}) {
	logrus.Error(l.generateMessage(args...)...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	logrus.Error(l.generateMessage(fmt.Sprintf(format, args...))...)
}

func (l *logger) Stopwatch(message string) func() {
	started := time.Now()
	l.Infof("STARTED: %s", message)

	return func() {
		l.Infof("DONE (%v): %s", time.Since(started), message)
	}
}

func (l *logger) Stopwatchf(format string, args ...interface{}) func() {
	return l.Stopwatch(fmt.Sprintf(format, args...))
}

func (l *logger) WithCtx(ctx context.Context) Logger {
	return withCtx(l.component, ctx)
}

type dummyLogger struct{}

func (l *dummyLogger) Info(v ...interface{})                             {}
func (l *dummyLogger) Infof(format string, v ...interface{})             {}
func (l *dummyLogger) Debug(v ...interface{})                            {}
func (l *dummyLogger) Debugf(format string, v ...interface{})            {}
func (l *dummyLogger) Warn(v ...interface{})                             {}
func (l *dummyLogger) Warnf(format string, v ...interface{})             {}
func (l *dummyLogger) Error(v ...interface{})                            {}
func (l *dummyLogger) Errorf(format string, v ...interface{})            {}
func (l *dummyLogger) Stopwatch(message string) func()                   { return nil }
func (l *dummyLogger) Stopwatchf(format string, v ...interface{}) func() { return nil }
func (l *dummyLogger) WithCtx(ctx context.Context) Logger                { return l }

func NonNullLogger(logger Logger) Logger {
	if logger != nil {
		return logger
	}
	return &dummyLogger{}
}

func init() {

}
