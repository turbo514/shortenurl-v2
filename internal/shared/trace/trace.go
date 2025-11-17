package mytrace

import (
	"context"
	"fmt"
	"github.com/turbo514/shortenurl-v2/shared/commonconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var serviceName string

func GetTracer() trace.Tracer {
	return otel.Tracer(serviceName)
}

//type MySampler struct{}
//
//func (MySampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
//	if rand.Float64() < 0.1 {
//		return sdktrace.SamplingResult{
//			Decision:   sdktrace.RecordAndSample,
//			Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
//		}
//	} else {
//		return sdktrace.SamplingResult{
//			Decision:   sdktrace.Drop,
//			Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
//		}
//	}
//}
//
//func (MySampler) Description() string {
//	return "MySampler"
//}

func InitOpenTelemetry(ctx context.Context, jaegerConfig *commonconfig.JaegerConfig, serviceInfo *commonconfig.ServiceInfo) (*sdktrace.TracerProvider, error) {
	// 1. 配置 OTLP Exporter 连接到 Jaeger OTLP 端口
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", jaegerConfig.Host, jaegerConfig.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to Jaeger: %w", err)
	}

	// 2. 配置 Exporter (将追踪数据发送到哪里)
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// 3. 设置资源属性
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceInfo.Name), // 您的服务名称
			semconv.ServiceVersion(serviceInfo.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	serviceName = serviceInfo.Name

	// 4. 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // 批量发送数据
		sdktrace.WithResource(res),     // 关联资源属性
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.05))), // 采样策略
	)

	// 5. 将 TracerProvider 注册为全局默认提供者
	otel.SetTracerProvider(tp)

	// 6. 设置 Context 传播器，确保 trace-id 可以在服务间正确传递
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}
