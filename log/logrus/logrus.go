package logrus

import (
	"github.com/sirupsen/logrus"
	"github.com/yeezyi/lpkg/log"
	"time"
)

type Logger struct {
	logger *logrus.Entry
}

func NewLogger(opts ...Option) *Logger {
	var std = logrus.New()
	std.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
		FieldMap:        nil,
	})
	logger := &Logger{logger: logrus.NewEntry(std)}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

func (l *Logger) Log(level log.Level, args ...interface{}) {
	switch level {
	case log.PanicLevel:
		l.logger.Panic(args...)
	case log.FatalLevel:
		l.logger.Fatal(args...)
	case log.ErrorLevel:
		l.logger.Error(args...)
	case log.WarnLevel:
		l.logger.Warn(args...)
	case log.InfoLevel:
		l.logger.Info(args...)
	case log.DebugLevel:
		l.logger.Debug(args...)
	}
}

func (l *Logger) Logf(level log.Level, format string, args ...interface{}) {
	switch level {
	case log.PanicLevel:
		l.logger.Panicf(format, args...)
	case log.FatalLevel:
		l.logger.Fatalf(format, args...)
	case log.ErrorLevel:
		l.logger.Errorf(format, args...)
	case log.WarnLevel:
		l.logger.Warnf(format, args...)
	case log.InfoLevel:
		l.logger.Infof(format, args...)
	case log.DebugLevel:
		l.logger.Debugf(format, args...)
	}
}

func (l *Logger) With(key string, val interface{}) log.Logger {
	return &Logger{logger: l.logger.WithField(key, val)}
}
