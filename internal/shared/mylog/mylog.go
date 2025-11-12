package mylog

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
)

// 设置基础日志器，使用文本格式输出到 os.Stderr
var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

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
