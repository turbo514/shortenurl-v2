package command

import (
	"context"
	"fmt"
	"github.com/turbo514/shortenurl-v2/analytics/domain"
)

type CreateClickEventCommand struct {
	Events []*domain.ClickEvent
}

type CreateClickEventHandler struct {
	repo         domain.IRepository
	clickCounter domain.IClickCounter
}

func NewCreateClickEventHandler(repo domain.IRepository, clickcounter domain.IClickCounter) *CreateClickEventHandler {
	return &CreateClickEventHandler{
		repo:         repo,
		clickCounter: clickcounter,
	}
}

// Handle 向写优化事实表插入点击记录,增加点击量
func (h *CreateClickEventHandler) Handle(ctx context.Context, cmd CreateClickEventCommand) error {
	if err := h.repo.InsertClickEvents(ctx, cmd.Events); err != nil {
		return fmt.Errorf("插入点击事件失败: %w", err)
	}

	if err := h.clickCounter.IncreaseMany(ctx, cmd.Events); err != nil {
		return fmt.Errorf("增加点击量失败: %w", err)
	}

	return nil
}
