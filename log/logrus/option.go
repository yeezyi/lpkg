package logrus

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Option func(logger *Logger)

func WithFormatter(formatterType string, logPretty bool) Option {
	return func(l *Logger) {
		if formatterType == "json" {
			l.logger.Logger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.DateTime,
				PrettyPrint:     logPretty,
			})
		} else if formatterType == "text" {
			l.logger.Logger.SetFormatter(&logrus.TextFormatter{
				TimestampFormat: time.DateTime,
				FullTimestamp:   true,
				PadLevelText:    true,
			})
		}
	}
}

func WithEnableCaller(caller bool) Option {
	return func(l *Logger) {
		l.logger.Logger.SetReportCaller(caller)
	}
}

func WithLevel(level logrus.Level) Option {
	return func(l *Logger) {
		l.logger.Logger.SetLevel(level)
	}
}
