package rabbitmq_publisher

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/turbo514/shortenurl-v2/link/adapter"
	"github.com/turbo514/shortenurl-v2/shared/dto"
	"github.com/turbo514/shortenurl-v2/shared/rabbitmq"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
)

var _ adapter.IEventPublisher = (*EventPublisher)(nil)

type EventPublisher struct {
	// TODO: 换成channel池
	conn    *amqp.Connection
	closeCh chan struct{}
}

func NewEventPublisher(conn *amqp.Connection) *EventPublisher {
	return &EventPublisher{
		conn:    conn,
		closeCh: make(chan struct{}),
	}
}

func (e *EventPublisher) Init() error {
	channel, err := e.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建初始化channel失败: %w", err)
	}

	if err := rabbitmq.InitClickEventExchange(channel); err != nil {
		return fmt.Errorf("初始化EventPublisher失败: %w", err)
	}
	if err := rabbitmq.InitCreateLinkEventExchange(channel); err != nil {
		return fmt.Errorf("初始化EventPublisher失败: %w", err)
	}
	if err := rabbitmq.InitClickEventQueue(channel); err != nil {
		return fmt.Errorf("初始化EventPublisher失败: %w", err)
	}
	if err := rabbitmq.BindAnalyticsClickEventQueue(channel); err != nil {
		return fmt.Errorf("初始化EventPublisher失败: %w", err)
	}
	if err := rabbitmq.InitCreateShortLinkQueue(channel); err != nil {
		return fmt.Errorf("初始化EventPublisher失败: %w", err)
	}
	if err := rabbitmq.BindUpdateCreateLinkEventQueue(channel); err != nil {
		return fmt.Errorf("初始化EventPublisher失败: %w", err)
	}

	return nil
}

func (e *EventPublisher) PublishClickEvent(ctx context.Context, event dto.ClickEvent) error {
	tr := mytrace.GetTracer()
	ctx, span := tr.Start(ctx, "EventPublisher.PublishClickEvent")
	defer span.End()

	channel, err := e.conn.Channel()
	defer channel.Close()

	if err != nil {
		return fmt.Errorf("创建ClickEventChannel失败: %w", err)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json序列化失败: %w", err)
	}

	err = channel.PublishWithContext(ctx, rabbitmq.ClickEventExchange, rabbitmq.ClickLink, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("发送到消息队列失败: %w", err)
	}

	return nil
}

func (e *EventPublisher) PublishCreateEvent(ctx context.Context, event dto.CreateEvent) error {
	channel, err := e.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建ClickEventChannel失败: %w", err)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json序列化失败: %w", err)
	}

	err = channel.Publish(rabbitmq.CreateLinkEventExchange, rabbitmq.CreateLink, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("发送到消息丢列失败: %w", err)
	}

	return nil
}

func (e *EventPublisher) Close() {
	close(e.closeCh)
}
