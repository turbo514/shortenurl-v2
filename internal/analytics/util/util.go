package util

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/turbo514/shortenurl-v2/shared/rabbitmq"
)

func InitQueue(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("create channel err: %v", err)
	}
	defer ch.Close()
	if err := rabbitmq.InitClickEventQueue(ch); err != nil {
		return err
	}
	if err := rabbitmq.InitCreateShortLinkQueue(ch); err != nil {
		return err
	}
	return nil
}
