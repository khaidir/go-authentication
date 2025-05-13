package logger

import (
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	Log        *zap.Logger
	once       sync.Once
	TraceIDKey = "trace_id"
)

func InitLogger() {
	once.Do(func() {
		var err error
		Log, err = zap.NewProduction()
		if err != nil {
			log.Fatal("failed to initialize logger: ", err)
		}
	})
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

// Logger defines the interface for logging.
type Logger interface {
	Info(msg string, err error, traceID string)
	Warn(msg string, err error, traceID string)
	Error(msg string, err error, traceID string)
	Debug(msg string, err error, traceID string)
}

// FileLogger implements Logger interface for file-based logging.
type FileLogger struct {
	logger *zap.Logger
}

func NewFileLogger(filePath string) (*FileLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{filePath}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &FileLogger{logger: logger}, nil
}

func (f *FileLogger) Info(msg string, err error, traceID string) {
	f.logger.Info(msg, zap.Error(err), zap.String("trace_id", traceID))
}

func (f *FileLogger) Warn(msg string, err error, traceID string) {
	f.logger.Warn(msg, zap.Error(err), zap.String("trace_id", traceID))
}

func (f *FileLogger) Error(msg string, err error, traceID string) {
	f.logger.Error(msg, zap.Error(err), zap.String("trace_id", traceID))
}

func (f *FileLogger) Debug(msg string, err error, traceID string) {
	f.logger.Debug(msg, zap.Error(err), zap.String("trace_id", traceID))
}

// LogWithTrace handles logging with trace ID using a logger that implements the Logger interface.
func LogWithTrace(c *gin.Context, msg string, err error, level string) {
	traceID, _ := c.Get(TraceIDKey)

	// Ensure logs directory exists
	if err := ensureLogsDirectory(); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		return
	}

	// Map log levels to file paths
	logFiles := map[string]string{
		"info":  "storage/logs/info.log",
		"warn":  "storage/logs/warn.log",
		"error": "storage/logs/error.log",
		"debug": "storage/logs/debug.log",
	}

	// Get the file path for the current log level
	filePath, ok := logFiles[level]
	if !ok {
		filePath = logFiles["info"] // Default to info.log if level is invalid
	}

	// Initialize logger
	fileLogger, logErr := NewFileLogger(filePath)
	if logErr != nil {
		log.Printf("Failed to create file logger: %v", logErr)
		return
	}
	defer fileLogger.logger.Sync()

	// Log the message based on the level
	switch level {
	case "info":
		fileLogger.Info(msg, err, traceID.(string))
	case "warn":
		fileLogger.Warn(msg, err, traceID.(string))
	case "error":
		fileLogger.Error(msg, err, traceID.(string))
	case "debug":
		fileLogger.Debug(msg, err, traceID.(string))
	default:
		fileLogger.Info(msg, err, traceID.(string))
	}
}

// ensureLogsDirectory ensures the logs directory exists.
func ensureLogsDirectory() error {
	return os.MkdirAll("storage/logs", 0755)
}
