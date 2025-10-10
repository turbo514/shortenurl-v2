package myprom

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
)

// ExemplarFromContext 样本提取函数：从 Context 中提取 Trace ID 作为 Exemplar
// 用于将指标 (Metrics) 和追踪 (Traces) 关联起来。
func ExemplarFromContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{"traceID": span.TraceID().String()}
	}
	return nil
}

// LabelsFromContext 动态标签提取函数：从 gRPC 元数据 (Metadata) 中提取租户ID ('tenant-id') 作为 Prometheus 标签。
func LabelsFromContext(ctx context.Context) prometheus.Labels {
	labels := prometheus.Labels{}

	// 从 Context 中提取传入的 gRPC 元数据
	md := metadata.ExtractIncoming(ctx)
	// 获取 'tenant-name' 字段的值
	tenantID := md.Get("tenant-id")
	if tenantID == "" {
		tenantID = "unknown"
	}
	labels["tenant_id"] = tenantID

	return labels
}
