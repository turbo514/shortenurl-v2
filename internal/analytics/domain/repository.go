package domain

import "context"

type IRepository interface {
	CreateLinks(ctx context.Context, links []*Link) error
	InsertClickEvents(ctx context.Context, clickEvents []*ClickEvent) error
}
