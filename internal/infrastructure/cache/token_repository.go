package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	ch *redis.Client
}

func NewTokenRepository(ch *redis.Client) repository.TokenRepository {
	return &tokenRepository{ch}
}

func (*tokenRepository) key(userID uuid.UUID) string {
	return fmt.Sprintf("refresh:%s", userID.String())
}

func (r *tokenRepository) Get(ctx context.Context, userID uuid.UUID) (string, error) {
	v, err := r.ch.Get(ctx, r.key(userID)).Result()
	if err != nil {
		return "", err
	}

	return v, nil
}

func (r *tokenRepository) Set(ctx context.Context, token *domain.Token) error {
	err := r.ch.Set(ctx, r.key(token.UserID), token.Refresh, time.Duration(token.ExpirationHours*float32(time.Hour))).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	err := r.ch.Del(ctx, r.key(userID)).Err()
	if err != nil {
		return err
	}

	return nil
}
