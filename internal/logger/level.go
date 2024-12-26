package logger

import "log/slog"

type LoggerLevel string

const (
	Info  LoggerLevel = "info"
	Warn  LoggerLevel = "warn"
	Error LoggerLevel = "error"
	Debug LoggerLevel = "debug"
)

func Into(l string) LoggerLevel {
	return LoggerLevel(l)
}

func (l LoggerLevel) ToSlogLevel() slog.Level {
	switch l {
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	case Debug:
		return slog.LevelDebug
	}

	return slog.LevelInfo
}
