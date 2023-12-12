package application

import (
	"context"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type UserService interface {
	Login(ctx context.Context, user *domain.User) (*dto.Tokens, error)
	Logout(ctx context.Context, userID, channelID uuid.UUID)
	GetTokens(ctx context.Context, userID uuid.UUID) (*dto.Tokens, error)
	Get(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	PreCreate(ctx context.Context, user *domain.CreateUser) (uuid.UUID, error)
	Create(ctx context.Context, preUserID uuid.UUID) (uuid.UUID, error)
	Update(ctx context.Context, user *domain.UpdateUser) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Verify(ctx context.Context, userID uuid.UUID, password string) error
	SendConfirmToken(ctx context.Context, email string) error
	SendResetToken(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, subID uuid.UUID, password string) error
}

type userService struct {
	c   *config.Config
	pur repository.PreUserRepository
	ur  repository.UserRepository
	uar repository.UserActionRepository
	as  RESTAuthService
	ms  MailService
	cs  CryptoService
	es  EventService
}

func NewUserService(
	c *config.Config,
	pur repository.PreUserRepository,
	ur repository.UserRepository,
	uar repository.UserActionRepository,
	as RESTAuthService,
	ms MailService,
	cs CryptoService,
	es EventService,
) UserService {
	return &userService{c, pur, ur, uar, as, ms, cs, es}
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

	err = u.cs.CompareHash(found.Password, user.Password)
	if err != nil {
		return nil, ae.ErrInvalidPassword
	}

	tokens, err := u.as.GenerateTokens(ctx, found)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (u *userService) Logout(ctx context.Context, userID, channelID uuid.UUID) {
	if channelID == uuid.Nil {
		return
	}

	u.es.Publish(ctx, domain.Event{
		ID: uuid.New(),
		Data: domain.EventData{
			Action: enum.ChannelClosedAction,
		},
	}, userID, channelID)
}

func (u *userService) GetTokens(ctx context.Context, userID uuid.UUID) (*dto.Tokens, error) {
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

	hash, err := u.cs.GenerateHash(user.Password)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "userService.PreCreate.GenerateHash")
	}
	user.Password = hash

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
		hash, err := u.cs.GenerateHash(user.Password.String)
		if err != nil {
			return errors.Wrap(err, "userService.Update.GenerateHash")
		}
		user.Password = t.NewString(hash, false, true)
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

	err = u.cs.CompareHash(user.Password, password)
	if err != nil {
		return ae.ErrInvalidPassword
	}

	return nil
}

func (u *userService) SendConfirmToken(ctx context.Context, email string) error {
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

func (u *userService) SendResetToken(ctx context.Context, email string) (string, error) {
	subID := uuid.New()
	hashSubID, err := u.cs.GenerateHash(subID.String())
	if err != nil {
		return "", err
	}

	user, err := u.ur.GetByEmail(ctx, email)
	if err != nil {
		return hashSubID, err
	}

	token, err := u.as.GenerateResetToken(ctx, subID)
	if err != nil {
		return "", err
	}

	dur := time.Duration(u.c.JWT.ResetLifespanMinutes * float32(time.Minute))
	u.uar.Set(ctx, enum.UserResetPassword, subID, user.ID, dur)

	go u.ms.SendPasswordReset(email, token)

	return hashSubID, nil
}

func (u *userService) ResetPassword(ctx context.Context, subID uuid.UUID, password string) error {
	defer u.uar.Delete(ctx, enum.UserResetPassword, subID)
	userID, err := u.uar.Get(ctx, enum.UserResetPassword, subID)
	if err != nil {
		return err
	}

	err = u.Update(ctx, &domain.UpdateUser{
		ID:       userID,
		Password: t.NewString(password, false, true),
	})
	if err != nil {
		return err
	}

	err = u.as.RevokeRefreshTokens(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
