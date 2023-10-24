package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/google/uuid"
)

type PreUserRepository interface {
	Set(ctx context.Context, preUser *domain.CreateUser, expirationHours float32) (uuid.UUID, error)
	Get(ctx context.Context, preUserID uuid.UUID) (*domain.CreateUser, error)
	Delete(ctx context.Context, preUserID uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (uuid.UUID, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
