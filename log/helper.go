package log

type Helper struct {
	log Logger
}

func NewHelper(log Logger) *Helper {
	return &Helper{log: log}
}

func (h *Helper) Debugf(format string, args ...interface{}) { h.log.Logf(DebugLevel, format, args...) }
func (h *Helper) Infof(format string, args ...interface{})  { h.log.Logf(InfoLevel, format, args...) }
func (h *Helper) Warnf(format string, args ...interface{})  { h.log.Logf(WarnLevel, format, args...) }
func (h *Helper) Errorf(format string, args ...interface{}) { h.log.Logf(ErrorLevel, format, args...) }
func (h *Helper) Fatalf(format string, args ...interface{}) { h.log.Logf(FatalLevel, format, args...) }
func (h *Helper) Panicf(format string, args ...interface{}) { h.log.Logf(PanicLevel, format, args...) }

func (h *Helper) Debug(args ...interface{}) { h.log.Log(DebugLevel, args...) }
func (h *Helper) Info(args ...interface{})  { h.log.Log(InfoLevel, args...) }
func (h *Helper) Warn(args ...interface{})  { h.log.Log(WarnLevel, args...) }
func (h *Helper) Error(args ...interface{}) { h.log.Log(ErrorLevel, args...) }
func (h *Helper) Fatal(args ...interface{}) { h.log.Log(FatalLevel, args...) }
func (h *Helper) Panic(args ...interface{}) { h.log.Log(PanicLevel, args...) }

func (h *Helper) Printf(format string, args ...interface{}) {
	h.log.Logf(InfoLevel, format, args...)
}

func (h *Helper) With(key string, val interface{}) *Helper {
	return &Helper{h.log.With(key, val)}
}
