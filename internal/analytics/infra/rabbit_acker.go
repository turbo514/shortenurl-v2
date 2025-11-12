package infra

import "github.com/rabbitmq/amqp091-go"

type RabbitAcker struct {
	msg *amqp091.Delivery
}

func NewRabbitAcker(msg *amqp091.Delivery) *RabbitAcker {
	return &RabbitAcker{msg: msg}
}

func (r *RabbitAcker) Ack() error {
	return r.msg.Ack(true)
}

func (r *RabbitAcker) Nack() error {
	return r.msg.Nack(true, true)
}
