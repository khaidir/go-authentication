package logger

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var Log *zap.Logger

const TraceIDKey = "trace_id"

func InitLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set(TraceIDKey, traceID)
		c.Writer.Header().Set("X-Request-ID", traceID)
		c.Next()
	}
}
