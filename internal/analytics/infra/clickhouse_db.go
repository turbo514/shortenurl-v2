package infra

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
)

type ClickhouseDb struct {
	conn driver.Conn
}

var _ domain.IRepository = (*ClickhouseDb)(nil)

func NewClickhouseDb(conn driver.Conn) *ClickhouseDb {
	return &ClickhouseDb{conn: conn}
}

func (c *ClickhouseDb) CreateLinks(ctx context.Context, links []*domain.Link) error {
	//TODO implement me
	panic("implement me")
}

func (c *ClickhouseDb) InsertClickEvents(ctx context.Context, events []*domain.ClickEvent) error {
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO default.click_events_write")
	if err != nil {
		return fmt.Errorf("创建批处理失败: %w", err)
	}
	defer batch.Close()

	for _, event := range events {
		if err := batch.Append(
			event.ID,
			event.LinkID,
			event.ClickTime,
			event.ClickIP,
			event.UserAgent,
			event.Referrer,
		); err != nil {
			return fmt.Errorf("append失败: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("发送批处理失败: %w", err)
	}

	return nil
}
