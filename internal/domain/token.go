package domain

import (
	"github.com/google/uuid"
)

type Token struct {
	UserID          uuid.UUID
	Refresh         string
	ExpirationHours float32
}
