package infra

import (
	"context"
	"encoding/json"
	"github.com/turbo514/shortenurl-v2/analytics/controller"
	"github.com/turbo514/shortenurl-v2/analytics/cqrs/command"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
	"github.com/turbo514/shortenurl-v2/shared/dto"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	"github.com/twmb/franz-go/pkg/kgo"
	"sync"
)

type tp struct {
	topic     string
	partition int32
}

type KafkaClickEventReceiver struct {
	consumers map[tp]*controller.ClickEventHandler
	handler   *command.CreateClickEventHandler
}

func NewKafkaClickEventReceiver(handler *command.CreateClickEventHandler) *KafkaClickEventReceiver {
	return &KafkaClickEventReceiver{
		consumers: make(map[tp]*controller.ClickEventHandler),
		handler:   handler,
	}
}

func (r *KafkaClickEventReceiver) Poll(ctx context.Context, maxPollRecords int, client *kgo.Client) {
	logger := mylog.GetLogger()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fetches := client.PollRecords(ctx, maxPollRecords)

		if fetches.IsClientClosed() {
			return
		}

		fetches.EachError(func(_ string, _ int32, err error) {
			logger.Warn("KafkaClickEventReceiver.Poll Error", "err", err.Error())
		})

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			events := make([]*domain.ClickEvent, 0, len(p.Records))
			for i := range p.Records {
				var event dto.ClickEvent // FIXME: 对象池?
				if err := json.Unmarshal(p.Records[i].Value, &event); err != nil {
					// FIXME: 需要进行处理
					logger.Error("Failed to unmarshal click event",
						"error", err.Error(),
						"topic", p.Topic,
						"partition", p.Partition,
						"offset", p.Records[i].Offset)
				} else {
					events = append(events, &domain.ClickEvent{
						ID:        event.EventId,
						LinkID:    event.LinkID,
						TenantID:  event.TenantID,
						ClickTime: event.ClickTime,
						ClickIP:   event.IpAddress,
						UserAgent: event.UserAgent,
						Referrer:  event.Referrer,
					})
				}
			}

			if len(events) == 0 {
				return
			}
			tp := tp{p.Topic, p.Partition}

			// 由于 BlockRebalanceOnPoll，保证分区消费者存在
			r.consumers[tp].Recs <- controller.ClickEventWithAcker{ // // 将批次消息发送给对应 partition 的 goroutine
				Event: events,
				Acker: NewKafkaAcker(client, p.Records),
			}
		})

		client.AllowRebalance()
	}
}

// Assigned 分区被分配时回调
func (r *KafkaClickEventReceiver) Assigned(ctx context.Context, cl *kgo.Client, assigned map[string][]int32) {
	logger := mylog.GetLogger()

	for topic, partitions := range assigned {
		for _, partition := range partitions {
			c := controller.NewClickEventHandler(r.handler)
			r.consumers[tp{topic, partition}] = c
			logger.Debug("启动一个工作协程", "topic", topic, "partition", partition)
			go c.Consume() // 启动独立 goroutine 消费该分区
		}
	}
}

// Lost 分区丢失/撤销回调
func (r *KafkaClickEventReceiver) Lost(ctx context.Context, cl *kgo.Client, lost map[string][]int32) {
	logger := mylog.GetLogger()
	var wg sync.WaitGroup
	defer wg.Wait() // 等待所有分区消费者退出完成

	for topic, partitions := range lost {
		for _, partition := range partitions {
			tp := tp{topic, partition}
			c := r.consumers[tp]
			delete(r.consumers, tp)
			c.Quit() // 发送退出信号
			logger.Info("正在等待工作协程关闭", "topic", topic, "partition", partition)

			wg.Add(1)
			go func() { <-c.Done(); wg.Done() }() // 等待 goroutine 完成
		}
	}
}
