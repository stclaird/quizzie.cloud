package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes the global logger
func InitLogger() error {
	return InitLoggerWithConfig("", "development")
}

// InitLoggerWithLevel initializes the global logger with a specific log level (deprecated, use InitLoggerWithConfig)
func InitLoggerWithLevel(configLogLevel string) error {
	return InitLoggerWithConfig(configLogLevel, "development")
}

// InitLoggerWithConfig initializes the global logger with log level and environment from config
func InitLoggerWithConfig(configLogLevel, configEnvironment string) error {
	var config zap.Config

	// Get log level from config (default to InfoLevel)
	logLevel := getLogLevel(configLogLevel)

	// Use config environment (no fallback to gin mode detection)
	env := configEnvironment
	if env == "" {
		env = "development" // Default fallback
	}

	// Use development config for local development, production for deployment
	if env == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	// Replace the global logger
	zap.ReplaceGlobals(logger)
	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if err := zap.L().Sync(); err != nil {
		// In some environments (like tests), sync might fail
		// Log to stderr but don't panic
	}
}

// getLogLevel parses the log level from config and returns the appropriate zapcore.Level
// Defaults to InfoLevel if not specified or invalid
func getLogLevel(configLogLevel string) zapcore.Level {
	logLevelStr := strings.ToLower(configLogLevel)

	switch logLevelStr {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		// Default to InfoLevel if not specified or invalid
		return zap.InfoLevel
	}
}
