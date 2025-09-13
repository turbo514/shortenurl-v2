package entity

import (
	"time"

	"github.com/google/uuid"
)

type ShortLink struct {
	Code        string
	OriginalURL string
	TenantID    uuid.UUID
	ExpireAt    time.Time
}

func (l *ShortLink) IsExpired() bool {
	return time.Now().After(l.ExpireAt)
}
