package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
