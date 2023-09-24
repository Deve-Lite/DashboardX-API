package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/google/uuid"
)

type TokenRepository interface {
	SetRefresh(ctx context.Context, token *domain.Token) error
	GetRefresh(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteRefresh(ctx context.Context, userID uuid.UUID) error
}
