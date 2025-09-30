package rabbitmq

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

const (
	ClickEventExchange      string = "su.ex.prod.clickevent"
	CreateLinkEventExchange string = "su.ex.prod.createlinkevent"
)

const (
	AnalyticsClickEventQueue   string = "su.q.prod.analytics.clickevent"
	UpdateCreateShortLinkQueue string = "su.q.prod.Update.createlink"
)

const (
	ClickLink  string = "shortlink.click"
	CreateLink string = "shortlink.create"
)

func InitClickEventExchange(channel *amqp091.Channel) error {
	if err := channel.ExchangeDeclare(ClickEventExchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("初始化ClickEventExchange失败: %w", err)
	}
	return nil
}

func InitCreateLinkEventExchange(channel *amqp091.Channel) error {
	if err := channel.ExchangeDeclare(CreateLinkEventExchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("初始化CreateLinkEventExchange失败: %w", err)
	}
	return nil
}

func InitAnalyticsClickEventQueue(channel *amqp091.Channel) error {
	if _, err := channel.QueueDeclare(AnalyticsClickEventQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("初始化AnalyticsClickEventQueue失败: %w", err)
	}

	return nil
}

func BindAnalyticsClickEventQueue(channel *amqp091.Channel) error {
	if err := channel.QueueBind(AnalyticsClickEventQueue, ClickLink, ClickEventExchange, false, nil); err != nil {
		return fmt.Errorf("绑定AnalyticsClickEventQueue失败: %w", err)
	}
	return nil
}

func InitUpdateCreateShortLinkQueue(channel *amqp091.Channel) error {
	if _, err := channel.QueueDeclare(UpdateCreateShortLinkQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("初始化UpdateCreateShortLinkQueue失败: %w", err)
	}
	return nil
}

func InitUpdateCreateLinkEventQueue(channel *amqp091.Channel) error {
	if err := channel.QueueBind(UpdateCreateShortLinkQueue, CreateLink, CreateLinkEventExchange, false, nil); err != nil {
		return fmt.Errorf("绑定UpdateCreateLinkEventQueue失败: %w", err)
	}
	return nil
}
