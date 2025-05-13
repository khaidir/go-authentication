package utils

import (
	"auth-services/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogError(c *gin.Context, code int, message string, err error) {
	traceID := c.GetString(logger.TraceIDKey)
	logger.Log.Error(message,
		zap.String("trace_id", traceID),
		zap.Error(err))
}
