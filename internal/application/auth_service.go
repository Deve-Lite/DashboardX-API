package application

import (
	"context"
	"log"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/domain/repository"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RESTAuthService interface {
	VerifyToken(ctx context.Context, token string, tokenType string) (*dto.RESTClaims, error)
	GenerateTokens(ctx context.Context, user *domain.User) (*dto.Tokens, error)
}

type restAuthService struct {
	c  *config.Config
	tr repository.TokenRepository
}

func NewRESTAuthService(c *config.Config, tr repository.TokenRepository) RESTAuthService {
	return &restAuthService{c, tr}
}

func (a *restAuthService) GenerateTokens(ctx context.Context, user *domain.User) (*dto.Tokens, error) {
	arc := dto.RESTClaims{
		IsAdmin: user.IsAdmin,
	}
	arc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(a.c.JWT.AccessLifespanHours * float32(time.Hour))))
	arc.Issuer = user.ID.String()
	arc.ID = uuid.NewString()

	rrc := dto.RESTClaims{
		IsAdmin: user.IsAdmin,
	}
	rrc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(a.c.JWT.RefreshLifespanHours * float32(time.Hour))))
	rrc.Issuer = user.ID.String()
	rrc.ID = uuid.NewString()

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

	if err := a.tr.SetRefresh(ctx, &domain.Token{
		UserID:          user.ID,
		Refresh:         rt,
		ExpirationHours: a.c.JWT.RefreshLifespanHours,
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
		found, err := a.tr.GetRefresh(ctx, uuid.MustParse(claims.Issuer))
		if err != nil {
			return nil, err
		}

		if found != token {
			return nil, ae.ErrInvalidRefreshToken
		}

		err = a.tr.DeleteRefresh(ctx, uuid.MustParse(claims.Issuer))
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}
