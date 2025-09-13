package entity

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type ClickEvent struct {
	OriginalURL string
	ClickTime   time.Time
	IpAddress   net.IP
	TenantID    uuid.UUID
}
