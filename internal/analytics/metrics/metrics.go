package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	MessageTypeClick = "click_event"
)

var (
	// MessageConsumedTotal 消费消息总数 (成功/失败)
	MessageConsumedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "message_consumed_total",
			Help: "消费到的消息总数 (按topic和status区分)",
		},
		[]string{"topic", "status"},
	)

	// MessageConsumeDuration 消息处理耗时
	MessageConsumeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "message_consume_duration_seconds",
			Help:    "消息消费处理耗时分布(按topic区分)",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	// ProcessingMessages 当前正在处理的消息数
	ProcessingMessages = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "message_processing_inflight",
			Help: "当前正在处理的消息数量",
		},
	)

	// MqMessagesReceivedTotal 各类消息接收数量
	// 监控入口流量
	MqMessagesReceivedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mq_messages_received_total",
			Help: "各类消息接收数量",
		},
		[]string{"topic"},
	)

	// EventBatchSize 各类事件批处理大小分布
	EventBatchSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "event_batch_size",
			Help: "各类事件批处理大小分布",
		},
		[]string{"topic"},
	)
)
