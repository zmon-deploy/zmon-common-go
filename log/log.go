package log

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type LoggerFactory interface {
	GetLogger(component string) Logger
}

type loggerFactory struct {
	sentryWrapper SentryWrapper
}

func NewLoggerFactory(sentryDSN, ver string) (LoggerFactory, error) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	sw, err := newSentryWrapper(sentryDSN, ver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sentry wrapper")
	}

	return &loggerFactory{
		sentryWrapper: sw,
	}, nil
}

func (f *loggerFactory) GetLogger(component string) Logger {
	return newLogger(component, f.sentryWrapper)
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
}

type logger struct {
	component     string
	sentryWrapper SentryWrapper
}

func newLogger(component string, sentryWrapper SentryWrapper) Logger {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	return &logger{
		component:     component,
		sentryWrapper: sentryWrapper,
	}
}

func (l *logger) concatArgsWithComponent(args []interface{}) []interface{} {
	var arr []interface{}
	arr = append(arr, fmt.Sprintf("[%s] ", l.component))
	arr = append(arr, args...)
	return arr
}

func (l *logger) mergeFormatWithComponent(format string) string {
	return fmt.Sprintf("[%s] %s", l.component, format)
}

func (l *logger) Info(args ...interface{}) {
	logrus.Info(l.concatArgsWithComponent(args)...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	logrus.Infof(l.mergeFormatWithComponent(format), args...)
}

func (l *logger) Debug(args ...interface{}) {
	logrus.Debug(l.concatArgsWithComponent(args)...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	logrus.Debugf(l.mergeFormatWithComponent(format), args...)
}

func (l *logger) Warn(args ...interface{}) {
	logrus.Warn(l.concatArgsWithComponent(args)...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	logrus.Warnf(l.mergeFormatWithComponent(format), args...)
}

func (l *logger) Error(args ...interface{}) {
	arr := l.concatArgsWithComponent(args)
	logrus.Error(arr...)

	msg := fmt.Sprint(arr...)
	l.sentryWrapper.CaptureException(errors.New(msg))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(l.mergeFormatWithComponent(format), args...)
	logrus.Error(msg)

	l.sentryWrapper.CaptureException(errors.New(msg))
}

type dummyLogger struct{}

func (l *dummyLogger) Info(v ...interface{})                  {}
func (l *dummyLogger) Infof(format string, v ...interface{})  {}
func (l *dummyLogger) Debug(v ...interface{})                 {}
func (l *dummyLogger) Debugf(format string, v ...interface{}) {}
func (l *dummyLogger) Warn(v ...interface{})                  {}
func (l *dummyLogger) Warnf(format string, v ...interface{})  {}
func (l *dummyLogger) Error(v ...interface{})                 {}
func (l *dummyLogger) Errorf(format string, v ...interface{}) {}

func NonNullLogger(logger Logger) Logger {
	if logger != nil {
		return logger
	}
	return &dummyLogger{}
}
