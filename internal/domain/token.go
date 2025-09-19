package domain

import (
	"github.com/google/uuid"
	"net"
	"time"
)

type RefreshToken struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	HashedToken string
	ExpiresAt   time.Time
	Revoked     bool
	CreatedAt   time.Time
	LastUsedAt  time.Time
	DeviceInfo  string
	IPAddress   net.IP
}
