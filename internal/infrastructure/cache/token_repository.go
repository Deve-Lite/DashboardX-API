package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	ch *redis.Client
}

func NewTokenRepository(ch *redis.Client) repository.TokenRepository {
	return &tokenRepository{ch}
}

func (r *tokenRepository) GetRefresh(ctx context.Context, userID uuid.UUID) (string, error) {
	v, err := r.ch.Get(ctx, fmt.Sprintf("refresh/%s", userID)).Result()
	if err != nil {
		return "", err
	}

	return v, nil
}

func (r *tokenRepository) SetRefresh(ctx context.Context, token *domain.Token) error {
	err := r.ch.Set(ctx, fmt.Sprintf("refresh/%s", token.UserID), token.Refresh, time.Duration(token.ExpirationHours*float32(time.Hour))).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepository) DeleteRefresh(ctx context.Context, userID uuid.UUID) error {
	err := r.ch.Del(ctx, fmt.Sprintf("refresh/%s", userID)).Err()
	if err != nil {
		return err
	}

	return nil
}
