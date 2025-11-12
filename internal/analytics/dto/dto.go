package dto

import (
	"github.com/google/uuid"
)

type TopLinkView struct {
	ID          uuid.UUID
	OriginalURL string
	ClickTimes  int64
}

type TopLinks struct {
	List []TopLinkView
}
