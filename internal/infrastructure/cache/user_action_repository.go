package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type userActionRepository struct {
	ch *redis.Client
}

func NewUserActionRepository(ch *redis.Client) repository.UserActionRepository {
	return &userActionRepository{ch}
}

func (*userActionRepository) key(prefix enum.UserAction, keyID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", prefix, keyID.String())
}

func (r *userActionRepository) Get(ctx context.Context, prefix enum.UserAction, keyID uuid.UUID) (uuid.UUID, error) {
	v, err := r.ch.Get(ctx, r.key(prefix, keyID)).Result()
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.MustParse(v), nil
}

func (r *userActionRepository) Set(ctx context.Context, prefix enum.UserAction, keyID uuid.UUID, userID uuid.UUID, expiration time.Duration) error {
	err := r.ch.Set(ctx, r.key(prefix, keyID), userID.String(), expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *userActionRepository) Delete(ctx context.Context, prefix enum.UserAction, keyID uuid.UUID) error {
	err := r.ch.Del(ctx, r.key(prefix, keyID)).Err()
	if err != nil {
		return err
	}

	return nil
}
