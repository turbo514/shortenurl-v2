package events

import (
	"github.com/google/uuid"
	"net"
	"time"
)

type ClickEvent struct {
	LinkID    uuid.UUID `json:"link_id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ClickTime time.Time `json:"click_time"`
	IpAddress net.IP    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
}

type CreateEvent struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	LinkID      uuid.UUID `json:"link_id"`
	OriginalURL string    `json:"original_url"`
}
