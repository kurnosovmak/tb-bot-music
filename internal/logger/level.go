package logger

import (
	"errors"
	"log/slog"
	"strings"
)

// Константы уровней логирования
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// ErrLevelNotFound возвращается, если передан неизвестный уровень логирования
var ErrLevelNotFound = errors.New("unknown logger level")

// levelFromString конвертирует строку уровня в slog.Level
func levelFromString(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, ErrLevelNotFound
	}
}
