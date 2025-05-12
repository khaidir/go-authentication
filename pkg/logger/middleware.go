package logger

import (
	"auth-services/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		reqID, _ := c.Get(middleware.RequestIDKey)
		traceID, _ := c.Get(middleware.TraceIDKey)

		if Log != nil {
			Log.Info("HTTP Request",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.String("request_id", reqID.(string)),
				zap.String("trace_id", traceID.(string)),
				zap.Duration("duration", duration),
			)
		}
	}
}
