package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lmittmann/tint"
)

var defaultLogLevelMutex sync.RWMutex
var defaultLogLevel = slog.LevelInfo

// GetLogger returns a new logger
func GetLogger() *slog.Logger {
	defaultLogLevelMutex.RLock()
	defer defaultLogLevelMutex.RUnlock()
	return slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			TimeFormat: time.Kitchen,
			Level:      defaultLogLevel,
		}),
	)
}

// SetDefaultLevel sets the default log level
func SetDefaultLevel(level slog.Level) {
	defaultLogLevelMutex.Lock()
	defer defaultLogLevelMutex.Unlock()
	defaultLogLevel = level
}

// SetDefaultLevelFromString sets the default log level from a string
func SetDefaultLevelFromString(level string) error {
	switch strings.ToUpper(level) {
	case "DEBUG":
		SetDefaultLevel(slog.LevelDebug)
	case "INFO":
		SetDefaultLevel(slog.LevelInfo)
	case "WARN":
		SetDefaultLevel(slog.LevelWarn)
	case "ERROR":
		SetDefaultLevel(slog.LevelError)
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}
	return nil
}
