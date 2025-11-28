package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

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
	var err error

	// Use config environment (no fallback to environment variables)
	env := configEnvironment
	if env == "" {
		env = "development" // Default fallback
	}

	// Get log level from config (default to InfoLevel)
	logLevel := getLogLevel(configLogLevel)

	if env == "production" {
		// Production logger (JSON, structured)
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		Logger, err = config.Build()
	} else {
		// Development logger (console, human-readable)
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		Logger, err = config.Build()
	}

	if err != nil {
		return err
	}

	// Replace global logger
	zap.ReplaceGlobals(Logger)
	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if Logger == nil {
		// Fallback logger if not initialized
		Logger, _ = zap.NewDevelopment()
	}
	return Logger
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
