package mylog

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
)

var logLevel slog.Level

func init() {
	LogLevel := os.Getenv("LOG_LEVEL")
	switch LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

// 设置基础日志器，使用文本格式输出到 os.Stderr
var logger *slog.Logger

func SetLogLevel(level slog.Level) {
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

func GetLogger() *slog.Logger {
	return logger
}

// FillLogTraceID 从 Context 中提取 Trace ID 并作为日志字段
func FillLogTraceID(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}
