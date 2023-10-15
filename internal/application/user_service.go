package application

import (
	"context"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, user *domain.User) (*dto.Tokens, error)
	Refresh(ctx context.Context, userID uuid.UUID) (*dto.Tokens, error)
	Get(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error)
	Update(ctx context.Context, user *domain.UpdateUser) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Verify(ctx context.Context, userID uuid.UUID, password string) error
}

type userService struct {
	c  *config.Config
	ur repository.UserRepository
	a  RESTAuthService
}

func NewUserService(c *config.Config, ur repository.UserRepository, a RESTAuthService) UserService {
	return &userService{c, ur, a}
}

func (u *userService) Login(ctx context.Context, user *domain.User) (*dto.Tokens, error) {
	var found *domain.User
	var err error
	found, err = u.ur.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		return nil, ae.ErrInvalidPassword
	}

	tokens, err := u.a.GenerateTokens(ctx, found)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *userService) Refresh(ctx context.Context, userID uuid.UUID) (*dto.Tokens, error) {
	var user *domain.User
	var err error
	user, err = u.ur.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	tokens, err := u.a.GenerateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *userService) Get(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return u.ur.Get(ctx, userID)
}

func (u *userService) Create(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), int(u.c.Crytpo.HashCost))
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "userService.Create.GenerateFromPassword")
	}

	user.Password = string(hash)

	return u.ur.Create(ctx, user)
}

func (u *userService) Update(ctx context.Context, user *domain.UpdateUser) error {
	if user.Password.Set && !user.Password.Null {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password.String), int(u.c.Crytpo.HashCost))
		if err != nil {
			return errors.Wrap(err, "userService.Update.GenerateFromPassword")
		}

		user.Password = t.NewString(string(hash), false, true)
	}

	return u.ur.Update(ctx, user)
}

func (u *userService) Delete(ctx context.Context, userID uuid.UUID) error {
	return u.ur.Delete(ctx, userID)
}

func (u *userService) Verify(ctx context.Context, userID uuid.UUID, password string) error {
	var user *domain.User
	var err error
	user, err = u.ur.Get(ctx, userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return ae.ErrInvalidPassword
	}

	return nil
}
