package dto

import (
	"github.com/google/uuid"
	"net"
	"time"
)

// 消息队列传输用(服务间传输用)
type ClickEvent struct {
	EventId   uuid.UUID `json:"event_id"`
	LinkID    uuid.UUID `json:"link_id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ClickTime time.Time `json:"click_time"`
	IpAddress net.IP    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
}

type CreateEvent struct {
	EventId     uuid.UUID `json:"event_id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	LinkID      uuid.UUID `json:"link_id"`
	OriginalURL string    `json:"original_url"`
}
