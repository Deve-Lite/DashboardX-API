package cache

import (
	"context"
	"fmt"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
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

func (*tokenRepository) key(prefix enum.TokenType, ID, SubID uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", prefix, ID.String(), SubID.String())
}

func (*tokenRepository) keyAll(prefix enum.TokenType, SubID uuid.UUID) string {
	return fmt.Sprintf("%s:*:%s", prefix, SubID.String())
}

func (r *tokenRepository) Get(ctx context.Context, prefix enum.TokenType, ID, SubID uuid.UUID) (string, error) {
	v, err := r.ch.Get(ctx, r.key(prefix, ID, SubID)).Result()
	if err != nil {
		return "", err
	}

	return v, nil
}

func (r *tokenRepository) Set(ctx context.Context, token *domain.Token) error {
	err := r.ch.Set(ctx, r.key(token.Prefix, token.ID, token.SubID), token.Value, token.Expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepository) Delete(ctx context.Context, prefix enum.TokenType, ID, SubID uuid.UUID) error {
	err := r.ch.Del(ctx, r.key(prefix, ID, SubID)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepository) DeleteAll(ctx context.Context, prefix enum.TokenType, SubID uuid.UUID) error {
	keys, err := r.ch.Keys(ctx, r.keyAll(prefix, SubID)).Result()
	if err != nil {
		return err
	}

	for _, k := range keys {
		err := r.ch.Del(ctx, k).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
