package infra

import (
	"context"
	"github.com/turbo514/shortenurl-v2/analytics/controller"
	"github.com/twmb/franz-go/pkg/kgo"
)

var _ controller.IAcker = (*KafkaAcker)(nil)

// TODO: 尝试只考虑最后一个Record
type KafkaAcker struct {
	cl         *kgo.Client
	lastRecord []*kgo.Record
}

func NewKafkaAcker(cl *kgo.Client, lastRecord []*kgo.Record) *KafkaAcker {
	return &KafkaAcker{
		cl:         cl,
		lastRecord: lastRecord,
	}
}

func (k *KafkaAcker) Ack(ctx context.Context) error {
	return k.cl.CommitRecords(ctx, k.lastRecord...)
}

func (k *KafkaAcker) Nack(ctx context.Context) error {
	return nil
}
