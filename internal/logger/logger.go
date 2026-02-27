package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// NewLogger создаёт новый slog.Logger с JSON форматом и нужным уровнем
func NewLogger(rawLevel string) (*slog.Logger, error) {
	level, err := levelFromString(rawLevel)
	if err != nil {
		return nil, fmt.Errorf("error create logger: %w", &err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		// Кастомный формат времени: YYYY-MM-DD HH:MM:SS
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})

	return slog.New(handler), nil
}
