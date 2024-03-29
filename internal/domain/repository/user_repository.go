package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) bool
	Create(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error)
	Update(ctx context.Context, user *domain.UpdateUser) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
