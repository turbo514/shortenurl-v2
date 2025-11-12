package domain

import (
	"github.com/google/uuid"
	"time"
)

type Link struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	OriginalUrl string
	ShortCode   string
	UserId      uuid.UUID
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

func (l *Link) IsInvalid() bool {
	return l.ID == uuid.Nil
}

func (l *Link) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}
