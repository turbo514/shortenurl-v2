package query

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/turbo514/shortenurl-v2/analytics/dto"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	"gorm.io/gorm"
	"slices"
)

type GetTopLinksQuery struct {
	LinkToClickTimes map[uuid.UUID]int64
}

type GetTopLinksHandler struct {
	db    *gorm.DB
	cache *redis.Client
}

type TopLinkView struct {
	ID          uuid.UUID `gorm:"column:id;->;type:binary(16)"`
	OriginalURL string    `gorm:"column:original_url;->"`
}

func NewGetLinksHandler(db *gorm.DB, cache *redis.Client) *GetTopLinksHandler {
	return &GetTopLinksHandler{
		db:    db,
		cache: cache,
	}
}

func (h *GetTopLinksHandler) Handle(ctx context.Context, getLinksQuery *GetTopLinksQuery) (*dto.TopLinks, error) {
	// TODO: 从cache读取

	// 从排行榜读取目标id
	// FIXME: []byte 改成 uuid.UUID
	targetKeys := make([][]byte, 0, len(getLinksQuery.LinkToClickTimes))
	for k := range getLinksQuery.LinkToClickTimes {
		targetKeys = append(targetKeys, k[:])
	}

	// 从数据库获取链接详情
	links := make([]TopLinkView, 0)
	err := h.db.WithContext(ctx).Table("links").Select("id,original_url").Where("id IN ?", targetKeys).Find(&links).Error
	if err != nil {
		return nil, fmt.Errorf("query top links error: %w", err)
	}

	// 排序
	topLinks := dto.TopLinks{
		List: make([]dto.TopLinkView, len(links)),
	}
	for i := range links {
		topLinks.List[i].ClickTimes = getLinksQuery.LinkToClickTimes[links[i].ID]
		topLinks.List[i].OriginalURL = links[i].OriginalURL
		topLinks.List[i].ID = links[i].ID
	}

	slices.SortFunc(topLinks.List, func(a, b dto.TopLinkView) int { return int(b.ClickTimes - a.ClickTimes) })

	mylog.GetLogger().Debug("获取到排行榜链接详情", "排行榜", topLinks.List)

	// TODO: 排行榜存入cache中

	return &topLinks, nil
}
