package domain

import (
	"github.com/google/uuid"
	"net"
	"time"
)

type ClickEvent struct {
	ID        uuid.UUID
	LinkID    uuid.UUID
	TenantID  uuid.UUID
	ClickTime time.Time
	ClickIP   net.IP
	UserAgent string
	Referrer  string
}
