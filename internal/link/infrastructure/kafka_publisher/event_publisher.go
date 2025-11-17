package kafka_publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/shared/dto"
	"github.com/turbo514/shortenurl-v2/shared/keys"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
	"github.com/twmb/franz-go/pkg/kgo"
)

var _ adapter.IEventPublisher = (*KafkaEventPublisher)(nil)

type KafkaEventPublisher struct {
	client *kgo.Client
}

func NewKafkaEventPublisher(client *kgo.Client) *KafkaEventPublisher {
	return &KafkaEventPublisher{
		client: client,
	}
}

func (p *KafkaEventPublisher) PublishClickEvent(ctx context.Context, event dto.ClickEvent) error {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "KafkaEventPublisher.PublishClickEvent")
	defer span.End()

	v, err := json.Marshal(&event)
	if err != nil {
		return fmt.Errorf("json marshal event: %w", err)
	}

	record := kgo.Record{
		Topic: keys.ClickEventTopic,
		Value: v,
	}

	p.client.Produce(context.Background(), &record, func(_ *kgo.Record, err error) {
		if err != nil {
			mylog.GetLogger().Warn("点击事件发送失败", "err", err)
		}
	})
	return nil
}

func (p *KafkaEventPublisher) PublishCreateEvent(ctx context.Context, event dto.CreateEvent) error {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "KafkaEventPublisher.PublishCreateEvent")
	defer span.End()

	v, err := json.Marshal(&event)
	if err != nil {
		return fmt.Errorf("json marshal event: %w", err)
	}

	record := kgo.Record{
		Topic: keys.CreateEventTopic,
		Value: v,
	}

	p.client.Produce(context.Background(), &record, func(record *kgo.Record, err error) {
		if err != nil {
			mylog.GetLogger().Warn("创建事件发送失败", "err", err)
		}
	})

	return nil
}
