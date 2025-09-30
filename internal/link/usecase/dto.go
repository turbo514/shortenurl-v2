package usecase

import (
	"net"
	"time"
)

type ResolveRequest struct {
	Code      string
	ClickTime time.Time
	IpAddress net.IP
	UserAgent string
	Referrer  string
}
