package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
	"time"
)

type CreateLinkCommand struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	OriginalUrl string
	ShortUrl    string
	UserId      uuid.UUID
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type CreateLinkHandler struct {
	repo domain.IRepository
}

func (h *CreateLinkHandler) Handle(ctx context.Context, cmds []*CreateLinkCommand) error {
	links := make([]*domain.Link, len(cmds))
	for i, cmd := range cmds {
		link := &domain.Link{
			ID:          cmd.ID,
			TenantID:    cmd.TenantID,
			OriginalUrl: cmd.OriginalUrl,
			ShortCode:   cmd.ShortUrl,
			UserId:      cmd.UserId,
			CreatedAt:   cmd.CreatedAt,
			ExpiresAt:   cmd.ExpiresAt,
		}
		links[i] = link
	}
	
	return nil
}
