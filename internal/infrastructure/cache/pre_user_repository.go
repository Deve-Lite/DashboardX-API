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

type preUserRepository struct {
	ch *redis.Client
}

func NewPreUserRepository(ch *redis.Client) repository.PreUserRepository {
	return &preUserRepository{ch}
}

func (*preUserRepository) emailKey(email string) string {
	return fmt.Sprintf("pre-user:email:%s", email)
}

func (*preUserRepository) idKey(id uuid.UUID) string {
	return fmt.Sprintf("pre-user:id:%s", id.String())
}

func (r *preUserRepository) Get(ctx context.Context, preUserID uuid.UUID) (*domain.CreateUser, error) {
	var user domain.CreateUser
	if err := r.ch.HGetAll(ctx, r.idKey(preUserID)).Scan(&user); err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, redis.Nil
	}

	return &user, nil
}

func (r *preUserRepository) GetByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	preUserID, err := r.ch.Get(ctx, r.emailKey(email)).Result()
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.MustParse(preUserID), nil
}

func (r *preUserRepository) Set(ctx context.Context, preUser *domain.CreateUser, expirationHours float32) (uuid.UUID, error) {
	preUserID := uuid.New()

	exp := time.Duration(expirationHours * float32(time.Hour))

	pipe := r.ch.TxPipeline()

	pipe.HSet(ctx, r.idKey(preUserID), preUser)
	pipe.Expire(ctx, r.idKey(preUserID), exp)
	pipe.Set(ctx, r.emailKey(preUser.Email), preUserID.String(), exp)

	if _, err := pipe.Exec(ctx); err != nil {
		return uuid.Nil, err
	}

	return preUserID, nil
}

func (r *preUserRepository) Delete(ctx context.Context, preUserID uuid.UUID) error {
	user, err := r.Get(ctx, preUserID)
	if err != nil {
		return err
	}

	pipe := r.ch.TxPipeline()

	pipe.HDel(ctx, r.idKey(preUserID), "name", "password", "email", "is_admin", "language", "theme")
	pipe.Del(ctx, r.emailKey(user.Email))

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *preUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	v, err := r.ch.Exists(ctx, r.emailKey(email)).Result()
	if err != nil {
		return false, err
	}

	return v != 0, err
}
