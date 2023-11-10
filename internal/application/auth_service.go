package application

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/application/enum"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RESTAuthService interface {
	VerifyToken(ctx context.Context, token string, tokenType string) (*dto.RESTClaims, error)
	GenerateTokens(ctx context.Context, user *domain.User) (*dto.Tokens, error)
	VerifyConfirmToken(ctx context.Context, token string) (*jwt.RegisteredClaims, error)
	GenerateConfirmToken(ctx context.Context, preUserID uuid.UUID) (string, error)
	VerifyResetToken(ctx context.Context, token string) (*jwt.RegisteredClaims, error)
	GenerateResetToken(ctx context.Context, subID uuid.UUID) (string, error)
	VerifyResetPasswordSubject(ctx context.Context, subID uuid.UUID, enSubID string) error
	RevokeRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type restAuthService struct {
	c  *config.Config
	tr repository.TokenRepository
	cs CryptoService
}

func NewRESTAuthService(c *config.Config, tr repository.TokenRepository, cs CryptoService) RESTAuthService {
	return &restAuthService{c, tr, cs}
}

func (a *restAuthService) GenerateTokens(ctx context.Context, user *domain.User) (*dto.Tokens, error) {
	arc := dto.RESTClaims{
		IsAdmin: user.IsAdmin,
	}
	arc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(a.c.JWT.AccessLifespanHours * float32(time.Hour))))
	arc.Subject = user.ID.String()
	arc.ID = uuid.NewString()

	rrcDuration := time.Duration(a.c.JWT.RefreshLifespanHours * float32(time.Hour))
	rrc := dto.RESTClaims{
		IsAdmin: user.IsAdmin,
	}
	rrc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(rrcDuration))
	rrc.Subject = user.ID.String()
	rrcID := uuid.New()
	rrc.ID = rrcID.String()

	ac := jwt.NewWithClaims(jwt.SigningMethodHS256, arc)
	rc := jwt.NewWithClaims(jwt.SigningMethodHS256, rrc)

	var err error
	var at, rt string
	at, err = ac.SignedString([]byte(a.c.JWT.AccessSecret))
	if err != nil {
		return nil, err
	}
	rt, err = rc.SignedString([]byte(a.c.JWT.RefreshSecret))
	if err != nil {
		return nil, err
	}

	if err := a.tr.Set(ctx, &domain.Token{
		Prefix:     enum.TokenRefresh,
		ID:         rrcID,
		SubID:      user.ID,
		Value:      rt,
		Expiration: rrcDuration,
	}); err != nil {
		return nil, err
	}

	return &dto.Tokens{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (a *restAuthService) VerifyToken(ctx context.Context, token string, tokenType string) (*dto.RESTClaims, error) {
	var secret string
	if tokenType == "access" {
		secret = a.c.JWT.AccessSecret
	} else if tokenType == "refresh" {
		secret = a.c.JWT.RefreshSecret
	} else {
		log.Panic("invalid token type")
	}

	parsed, err := jwt.ParseWithClaims(token, &dto.RESTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := parsed.Claims.(*dto.RESTClaims)

	if tokenType == "refresh" {
		found, err := a.tr.Get(ctx, enum.TokenRefresh, uuid.MustParse(claims.ID), uuid.MustParse(claims.Subject))
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, ae.ErrTokenNotFound
			}
			return nil, err
		}

		if found != token {
			return nil, ae.ErrInvalidRefreshToken
		}

		err = a.tr.Delete(ctx, enum.TokenRefresh, uuid.MustParse(claims.ID), uuid.MustParse(claims.Subject))
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func (a *restAuthService) GenerateConfirmToken(ctx context.Context, preUserID uuid.UUID) (string, error) {
	rc := jwt.RegisteredClaims{}
	rc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(a.c.JWT.ConfirmLifespanHours * float32(time.Hour))))
	rc.Subject = preUserID.String()
	rc.ID = uuid.NewString()

	c := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)

	var err error
	var t string
	t, err = c.SignedString([]byte(a.c.JWT.ConfirmSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (a *restAuthService) VerifyConfirmToken(ctx context.Context, token string) (*jwt.RegisteredClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.c.JWT.ConfirmSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := parsed.Claims.(*jwt.RegisteredClaims)

	return claims, nil
}

func (a *restAuthService) GenerateResetToken(ctx context.Context, subID uuid.UUID) (string, error) {
	duration := time.Duration(a.c.JWT.ResetLifespanMinutes * float32(time.Minute))
	rc := jwt.RegisteredClaims{}
	rc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
	rc.Subject = subID.String()
	rcID := uuid.New()
	rc.ID = rcID.String()

	c := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)

	var err error
	var t string
	t, err = c.SignedString([]byte(a.c.JWT.ResetSecret))
	if err != nil {
		return "", err
	}

	a.tr.Set(ctx, &domain.Token{
		Prefix:     enum.TokenReset,
		ID:         rcID,
		SubID:      subID,
		Value:      t,
		Expiration: duration,
	})

	return t, nil
}

func (a *restAuthService) VerifyResetToken(ctx context.Context, token string) (*jwt.RegisteredClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.c.JWT.ResetSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := parsed.Claims.(*jwt.RegisteredClaims)

	found, err := a.tr.Get(ctx, enum.TokenReset, uuid.MustParse(claims.ID), uuid.MustParse(claims.Subject))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ae.ErrTokenNotFound
		}
		return nil, err
	}

	if found != token {
		return nil, ae.ErrInvalidRefreshToken
	}

	err = a.tr.Delete(ctx, enum.TokenReset, uuid.MustParse(claims.ID), uuid.MustParse(claims.Subject))
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (a *restAuthService) VerifyResetPasswordSubject(ctx context.Context, subID uuid.UUID, hashSubID string) error {
	if err := a.cs.CompareHash(hashSubID, subID.String()); err != nil {
		return err
	}

	return nil
}

func (a *restAuthService) RevokeRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return a.tr.DeleteAll(ctx, enum.TokenRefresh, userID)
}
