package repository

import (
	"context"
	"time"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/google/uuid"
)

type UserActionRepository interface {
	Set(ctx context.Context, prefix enum.UserAction, keyID uuid.UUID, userID uuid.UUID, expiration time.Duration) error
	Get(ctx context.Context, prefix enum.UserAction, keyID uuid.UUID) (uuid.UUID, error)
	Delete(ctx context.Context, prefix enum.UserAction, keyID uuid.UUID) error
}
