package log

type Logger interface {
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	With(key string, val interface{}) Logger
}

var log Logger

func SetLogger(logger Logger) {
	log = logger
}

func Default() Logger {
	return log
}

func Debugf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Logf(DebugLevel, format, args...)
}
func Infof(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Logf(InfoLevel, format, args...)
}
func Warnf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Logf(WarnLevel, format, args...)
}
func Errorf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Logf(ErrorLevel, format, args...)
}
func Fatalf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Logf(FatalLevel, format, args...)
}
func Panicf(format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Logf(PanicLevel, format, args...)
}

func Debug(args ...interface{}) {
	if log == nil {
		return
	}
	log.Log(DebugLevel, args...)
}
func Info(args ...interface{}) {
	if log == nil {
		return
	}
	log.Log(InfoLevel, args...)
}
func Warn(args ...interface{}) {
	if log == nil {
		return
	}
	log.Log(WarnLevel, args...)
}
func Error(args ...interface{}) {
	if log == nil {
		return
	}
	log.Log(ErrorLevel, args...)
}
func Fatal(args ...interface{}) {
	if log == nil {
		return
	}
	log.Log(FatalLevel, args...)
}
func Panic(args ...interface{}) {
	if log == nil {
		return
	}
	log.Log(PanicLevel, args...)
}

func With(key string, val interface{}) *Helper {
	if log == nil {
		return nil
	}
	return &Helper{log: log.With(key, val)}
}
