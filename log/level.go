package log

import "strings"

type Level int8

const (
	PanicLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func (l Level) String() string {
	switch l {
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

func StrToLevel(in string) Level {
	switch strings.ToUpper(in) {
	case "PANIC":
		return PanicLevel
	case "FATAL":
		return FatalLevel
	case "ERROR":
		return ErrorLevel
	case "WARN":
		return WarnLevel
	case "INFO":
		return InfoLevel
	case "DEBUG":
		return DebugLevel
	}
	return InfoLevel
}
