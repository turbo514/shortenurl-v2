package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/turbo514/shortenurl-v2/shared/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

type LinkMetrics struct {
	requestTotal    *prometheus.CounterVec   // 请求总数
	requestDuration *prometheus.HistogramVec // 请求用时直方图
}

// NewLinkMetrics 用于创建和初始化所有指标实例
func NewLinkMetrics() *LinkMetrics {
	return &LinkMetrics{
		requestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "link_service_requests_total",
				Help: "短链接服务的grpc请求总数",
			},
			[]string{"grpc_service", "grpc_method", "grpc_code"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "link_service_request_duration_seconds",
				Help:    "短链接服务的grpc请求响应用时直方图",
				Buckets: []float64{0.001, 0.01, 0.1, 1, 10}, // 小于1毫秒,小于10毫秒,小于100毫秒,小于1秒,小于10秒,,更大的
			},
			[]string{"grpc_service", "grpc_method"},
		),
	}
}

func (m *LinkMetrics) Describe(ch chan<- *prometheus.Desc) {
	// 调用内部指标的 Describe 方法，发送其 Descriptor
	m.requestTotal.Describe(ch)
	m.requestDuration.Describe(ch)
}
func (m *LinkMetrics) Collect(ch chan<- prometheus.Metric) {
	// 调用内部指标的 Collect 方法，发送它们当前的值快照
	m.requestTotal.Collect(ch)
	m.requestDuration.Collect(ch)
}

// UnaryServerInterceptor 返回一个 gRPC 拦截器函数, 负责在请求结束时调用 requestTotal.Inc() 和 requestDuration.Observe()
func (m *LinkMetrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 记录请求开始时间
		startTime := time.Now()

		// 调用下一个handler
		resp, err = handler(ctx, req)

		// 记录请求结束
		duration := time.Since(startTime)

		// 提取标签值
		serviceName, methodName := util.ParseFullMethod(info.FullMethod)
		statusCode := status.Code(err).String()

		// 记录请求总数
		m.requestTotal.WithLabelValues(serviceName, methodName, statusCode).Inc()
		// 记录请求延迟
		m.requestDuration.WithLabelValues(serviceName, methodName).Observe(duration.Seconds())

		return resp, err
	}
}
