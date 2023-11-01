package domain

import (
	"time"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/google/uuid"
)

type Token struct {
	Prefix     enum.TokenType
	ID         uuid.UUID
	SubID      uuid.UUID
	Value      string
	Expiration time.Duration
}
