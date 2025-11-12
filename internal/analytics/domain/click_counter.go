package domain

import (
	"context"
	"github.com/google/uuid"
)

type IClickCounter interface {
	//Increase(ctx context.Context, linkID uuid.UUID, increasement int) error
	IncreaseMany(ctx context.Context, links []*ClickEvent) error
	GetTopToday(ctx context.Context, num int64, tenant string) (map[uuid.UUID]int64, int64, error)
}
