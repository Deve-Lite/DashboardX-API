package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/google/uuid"
)

type TokenRepository interface {
	Set(ctx context.Context, token *domain.Token) error
	Get(ctx context.Context, userID uuid.UUID) (string, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}
