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
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, user *domain.User) (*dto.Tokens, error)
	Refresh(ctx context.Context, userID uuid.UUID) (*dto.Tokens, error)
	Get(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	PreCreate(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error)
	Create(ctx context.Context, preUserID uuid.UUID) (uuid.UUID, error)
	Update(ctx context.Context, user *domain.UpdateUser) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Verify(ctx context.Context, userID uuid.UUID, password string) error
	ResendConfirm(ctx context.Context, email string) error
}

type userService struct {
	c   *config.Config
	pur repository.PreUserRepository
	ur  repository.UserRepository
	as  RESTAuthService
	ms  MailService
}

func NewUserService(c *config.Config, pur repository.PreUserRepository, ur repository.UserRepository, as RESTAuthService, ms MailService) UserService {
	return &userService{c, pur, ur, as, ms}
}

func (u *userService) Login(ctx context.Context, user *domain.User) (*dto.Tokens, error) {
	var found *domain.User
	var err error
	found, err = u.ur.GetByEmail(ctx, user.Email)
	if err != nil {
		if b, _ := u.pur.ExistsByEmail(ctx, user.Email); b {
			return nil, ae.ErrConfirmationRequired
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		return nil, ae.ErrInvalidPassword
	}

	tokens, err := u.as.GenerateTokens(ctx, found)
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

	tokens, err := u.as.GenerateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *userService) Get(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return u.ur.Get(ctx, userID)
}

func (u *userService) PreCreate(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error) {
	if u.ur.ExistsByEmail(ctx, user.Email) {
		return uuid.Nil, ae.ErrEmailExists
	} else {
		if b, _ := u.pur.ExistsByEmail(ctx, user.Email); b {
			return uuid.Nil, ae.ErrEmailExists
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), int(u.c.Crytpo.HashCost))
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "userService.PreCreate.GenerateFromPassword")
	}

	user.Password = string(hash)

	preUserID, err := u.pur.Set(ctx, user, u.c.JWT.ConfirmLifespanHours)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "userService.PreCreate.SetPreUser")
	}

	token, err := u.as.GenerateConfirmToken(ctx, preUserID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "userService.PreCreate.GenerateConfirmToken")
	}

	go u.ms.SendConfirmAccount(user.Email, token)

	return preUserID, nil
}

func (u *userService) Create(ctx context.Context, preUserID uuid.UUID) (uuid.UUID, error) {
	preUser, err := u.pur.Get(ctx, preUserID)
	if err != nil {
		return uuid.Nil, ae.ErrNoAwaitingConfirm
	}

	userID, err := u.ur.Create(ctx, preUser)
	if err != nil {
		return uuid.Nil, ae.ErrUserCreation
	}

	if err := u.pur.Delete(ctx, preUserID); err != nil {
		go u.ur.Delete(ctx, userID)

		return uuid.Nil, ae.ErrUserCreation
	}

	return userID, nil
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

func (u *userService) ResendConfirm(ctx context.Context, email string) error {
	preUserID, err := u.pur.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ae.ErrNoAwaitingConfirm
		}
		return err
	}

	token, err := u.as.GenerateConfirmToken(ctx, preUserID)
	if err != nil {
		return errors.Wrap(err, "userService.ResendConfirm.GenerateConfirmToken")
	}

	go u.ms.SendConfirmAccount(email, token)

	return nil
}
