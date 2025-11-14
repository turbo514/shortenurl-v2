package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/turbo514/shortenurl-v2/shared/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

var metrics *MyMetrics = newLinkMetrics()

func GetListMetrics() *MyMetrics {
	return metrics

}

var _ prometheus.Collector = (*MyMetrics)(nil)

type MyMetrics struct {
	requestTotal *prometheus.CounterVec // gRPC 请求总数，标签：service, method, code

	distributedCacheLookUpTotal     *prometheus.CounterVec   // 分布式缓存命中情况,标签：result=hit/miss/error/hit_nil
	distributedCacheDurationSeconds *prometheus.HistogramVec // get/set 耗时,标签：operation=get/set

	localCacheLookUpTotal *prometheus.CounterVec // 本地缓存命中情况 // result=hit/miss/error

	dbDurationSeconds *prometheus.HistogramVec // 数据库操作耗时 // 标签：operation=query/insert/update
	dbOperationsTotal *prometheus.CounterVec   // 数据库操作总数 // 标签：operation=query/insert/update

	mqPublishTotal           *prometheus.CounterVec   // 消息队列发送情况 // 标签：event_type=click/create, result=success/fail
	mqPublishDurationSeconds *prometheus.HistogramVec // 标签：event_type=click/create
}

// NewLinkMetrics 用于创建和初始化所有指标实例
func newLinkMetrics() *MyMetrics {
	return &MyMetrics{
		requestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "link_service_requests_total",
				Help: "短链接服务的grpc请求总数",
			},
			[]string{"grpc_service", "grpc_method", "grpc_code"},
		),
		distributedCacheLookUpTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "link_service_distributed_cache_lookup_total",
				Help: "分布式缓存的命中情况",
			},
			[]string{"result"},
		),
		distributedCacheDurationSeconds: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "link_service_distributed_cache_duration_seconds",
				Help:    "分布式缓存的耗时直方图",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		localCacheLookUpTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "link_service_local_cache_lookup_total",
				Help: "本地缓存的命中/未命中数量",
			},
			[]string{"result"},
		),
		dbDurationSeconds: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "link_service_db_duration_seconds",
				Help:    "数据库的耗时直方图",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		dbOperationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "link_service_db_operations_total",
				Help: "数据库操作总数",
			},
			[]string{"operation"},
		),
		mqPublishTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "link_service_mq_publish_total",
				Help: "消息队列事件发送请求情况",
			},
			[]string{"event_type", "result"},
		),
		mqPublishDurationSeconds: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "link_service_mq_publish_duration_seconds",
				Help:    "消息队列发送消息的耗时直方图",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"event_type"},
		),
	}
}

func (m *MyMetrics) Describe(ch chan<- *prometheus.Desc) {
	// 调用内部指标的 Describe 方法，发送其 Descriptor
	m.requestTotal.Describe(ch)
	m.distributedCacheLookUpTotal.Describe(ch)
	m.distributedCacheDurationSeconds.Describe(ch)
	m.localCacheLookUpTotal.Describe(ch)
	m.dbDurationSeconds.Describe(ch)
	m.dbOperationsTotal.Describe(ch)
	m.mqPublishTotal.Describe(ch)
	m.mqPublishDurationSeconds.Describe(ch)
}
func (m *MyMetrics) Collect(ch chan<- prometheus.Metric) {
	// 调用内部指标的 Collect 方法，发送它们当前的值快照
	m.requestTotal.Collect(ch)
	m.distributedCacheLookUpTotal.Collect(ch)
	m.distributedCacheDurationSeconds.Collect(ch)
	m.localCacheLookUpTotal.Collect(ch)
	m.dbDurationSeconds.Collect(ch)
	m.dbOperationsTotal.Collect(ch)
	m.mqPublishTotal.Collect(ch)
	m.mqPublishDurationSeconds.Collect(ch)
}

// UnaryServerInterceptor 返回一个 gRPC 拦截器函数, 统计方法调用情况
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 调用下一个handler
		resp, err = handler(ctx, req)

		// 提取标签值
		serviceName, methodName := util.ParseFullMethod(info.FullMethod)
		statusCode := status.Code(err).String()

		// 记录请求总数
		metrics.requestTotal.WithLabelValues(serviceName, methodName, statusCode).Inc()

		return resp, err
	}
}

const (
	hits      = "hits"
	miss      = "miss"
	hitsNil   = "hits nil"
	hitsError = "error"
)

func AddDistributedCacheLookUpTotalHits() {
	metrics.distributedCacheLookUpTotal.WithLabelValues(hits).Inc()
}

func AddDistributedCacheLookUpTotalMiss() {
	metrics.distributedCacheLookUpTotal.WithLabelValues(miss).Inc()
}

func AddDistributedCacheLookUpTotalHitNil() {
	metrics.distributedCacheLookUpTotal.WithLabelValues(hitsNil).Inc()
}

func AddDistributedCacheLookUpTotalError() {
	metrics.distributedCacheLookUpTotal.WithLabelValues(hitsError).Inc()
}

func AddLocalCacheLookUpTotalHits() {
	metrics.localCacheLookUpTotal.WithLabelValues(hits).Inc()
}
func AddLocalCacheLookUpTotalMiss() {
	metrics.localCacheLookUpTotal.WithLabelValues(miss).Inc()
}
func AddLocalCacheLookUpTotalHitNil() {
	metrics.localCacheLookUpTotal.WithLabelValues(hitsNil).Inc()
}

const (
	get = "get"
	set = "set"
)

func ObserveDistributedCacheDurationSecondsGet(duration time.Duration) {
	metrics.distributedCacheDurationSeconds.WithLabelValues(get).Observe(duration.Seconds())
}
func ObserveDistributedCacheDurationSecondsSet(duration time.Duration) {
	metrics.distributedCacheDurationSeconds.WithLabelValues(set).Observe(duration.Seconds())
}

const (
	query  = "query"
	insert = "insert"
	update = "update"
)

func ObserveDbDurationSecondsQuery(duration time.Duration) {
	metrics.dbDurationSeconds.WithLabelValues(query).Observe(duration.Seconds())
}
func AddDbOperationsTotalQuery() {
	metrics.dbOperationsTotal.WithLabelValues(query).Inc()
}

func ObserveDbDurationSecondsInsert(duration time.Duration) {
	metrics.dbDurationSeconds.WithLabelValues(insert).Observe(duration.Seconds())
}
func AddDbOperationsTotalInsert() {
	metrics.dbOperationsTotal.WithLabelValues(insert).Inc()
}

const (
	click  = "click"
	create = "create"
)

func ObserveMqPublishDurationSecondsClick(duration time.Duration) {
	metrics.mqPublishDurationSeconds.WithLabelValues(click).Observe(duration.Seconds())
}
func ObserveMqPublishDurationSecondsCreate(duration time.Duration) {
	metrics.mqPublishDurationSeconds.WithLabelValues(create).Observe(duration.Seconds())
}

func AddMqPublishTotalClickSuccess() {
	metrics.mqPublishTotal.WithLabelValues(click, "success").Inc()
}
func AddMqPublishTotalClickFailed() {
	metrics.mqPublishTotal.WithLabelValues(click, "failed").Inc()
}
