package repository

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/google/uuid"
)

type TokenRepository interface {
	Set(ctx context.Context, token *domain.Token) error
	Get(ctx context.Context, prefix enum.TokenType, ID, SubID uuid.UUID) (string, error)
	Delete(ctx context.Context, prefix enum.TokenType, ID, SubID uuid.UUID) error
	DeleteAll(ctx context.Context, prefix enum.TokenType, SubID uuid.UUID) error
}
