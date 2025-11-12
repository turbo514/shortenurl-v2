package infra

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/turbo514/shortenurl-v2/analytics/controller"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
	dto "github.com/turbo514/shortenurl-v2/shared/dto"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
)

type RabbitMqClickEventReceiver struct {
	conn    *amqp.Connection
	queue   string
	channel chan controller.ClickEventWithAcker
	closeCh chan struct{}
}

func NewRabbitMqClickEventReceiver(conn *amqp.Connection, queue string) *RabbitMqClickEventReceiver {
	return &RabbitMqClickEventReceiver{
		conn:    conn,
		queue:   queue,
		channel: make(chan controller.ClickEventWithAcker),
		closeCh: make(chan struct{}),
	}
}

func (r *RabbitMqClickEventReceiver) GetChannel() chan controller.ClickEventWithAcker {
	return r.channel
}

func (r *RabbitMqClickEventReceiver) Close() {
	close(r.closeCh)
}

func (r *RabbitMqClickEventReceiver) Start() error {
	ch := r.channel

	channel, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建channel失败: %w", err)
	}
	msgs, err := channel.Consume(r.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("创建consumer失败: %w", err)
	}

	for {
		select {
		case <-r.closeCh:
			return nil
		case msg := <-msgs:
			body := msg.Body
			var dtoEvent dto.ClickEvent
			if err := json.Unmarshal(body, &dtoEvent); err != nil {
				mylog.GetLogger().Warn("消息解析失败", "err", err.Error())
				// FIXME 目前从简,不回归消息队列
				msg.Nack(false, false)
			} else {
				// debug
				mylog.GetLogger().Debug("监听到点击事件", "event", dtoEvent)
				ch <- controller.ClickEventWithAcker{
					Event: &domain.ClickEvent{
						ID:        dtoEvent.EventId,
						LinkID:    dtoEvent.LinkID,
						TenantID:  dtoEvent.TenantID,
						ClickTime: dtoEvent.ClickTime,
						ClickIP:   dtoEvent.IpAddress,
						UserAgent: dtoEvent.UserAgent,
						Referrer:  dtoEvent.Referrer,
					},
					Acker: NewRabbitAcker(&msg),
				}
			}
		}
	}
}
