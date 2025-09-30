package entity

import (
	"time"

	"github.com/google/uuid"
)

type ShortLink struct {
	ID          uuid.UUID
	ShortCode   string
	OriginalURL string
	TenantID    uuid.UUID
	UserID      uuid.UUID
	ExpireAt    *time.Time
}

func (l *ShortLink) IsExpired() bool {
	if l.ExpireAt != nil {
		return time.Now().After(*(l.ExpireAt))
	}
	return false
}

func (l *ShortLink) IsInvalid() bool {
	return l.ID == uuid.Nil
}
