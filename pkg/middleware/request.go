package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const (
	RequestIDKey = "X-Request-ID"
	TraceIDKey   = "Trace-ID"
)

func RequestContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(RequestIDKey)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(RequestIDKey, reqID)

		// Get traceID from OpenTelemetry context if exists
		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()
		c.Set(TraceIDKey, traceID)

		c.Writer.Header().Set(RequestIDKey, reqID)
		if traceID != "" {
			c.Writer.Header().Set(TraceIDKey, traceID)
		}

		c.Next()
	}
}
