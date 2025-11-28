package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZapMiddleware creates a Gin middleware using Zap for request logging
func GinZapMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the params
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		logger := zap.L().With(
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.Int("size", bodySize),
		)

		switch {
		case statusCode >= 500:
			logger.Error("Server error")
		case statusCode >= 400:
			logger.Warn("Client error")
		case statusCode >= 300:
			logger.Info("Redirection")
		default:
			logger.Info("Request completed")
		}
	}
}

// GinRecoveryWithZap creates a Gin middleware for panic recovery with Zap logging
func GinRecoveryWithZap(stack bool) gin.HandlerFunc {
	return gin.RecoveryWithWriter(&zapRecoveryWriter{stack: stack})
}

type zapRecoveryWriter struct {
	stack bool
}

func (w *zapRecoveryWriter) Write(p []byte) (n int, err error) {
	if w.stack {
		zap.L().Error("Panic recovered", zap.String("stack", string(p)))
	} else {
		zap.L().Error("Panic recovered", zap.String("message", string(p)))
	}
	return len(p), nil
}
