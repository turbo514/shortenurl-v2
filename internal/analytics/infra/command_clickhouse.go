package infra

//func (c *ClickhouseDb) InsertLinks(ctx context.Context, links []*domain.Link) error {
//	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO default.links")
//	if err != nil {
//		return fmt.Errorf("创建批处理失败: %w", err)
//	}
//	defer batch.Close()
//
//	for _, link := range links {
//		if err := batch.Append(
//			link.ID,
//			link.TenantID,
//			link.OriginalUrl,
//			link.ShortCode,
//			link.UserId,
//			link.CreatedAt,
//			link.ExpiresAt,
//		); err != nil {
//			return fmt.Errorf("append失败: %w", err)
//		}
//	}
//
//	if err := batch.Send(); err != nil {
//		return fmt.Errorf("发送批处理失败: %w", err)
//	}
//
//	return nil
//}
//
//func (c *ClickhouseDb) InsertClickEvent(ctx context.Context, events []*domain.ClickEvent) error {
//	panic("implement me")
//}
