package auth

import (
	"log"
	"time"

	"github.com/Deve-Lite/DashboardX-API-PoC/config"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type RESTAuth interface {
	VerifyToken(token string, tokenType string) (*RESTClaims, error)
	GenerateTokens(user *domain.User) (*Tokens, error)
}

type restAuth struct {
	c *config.Config
}

type RESTClaims struct {
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewRESTAuth(c *config.Config) RESTAuth {
	return &restAuth{c}
}

func (a *restAuth) GenerateTokens(user *domain.User) (*Tokens, error) {
	arc := RESTClaims{
		IsAdmin: user.IsAdmin,
	}
	arc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(a.c.JWT.AccessLifespan * float32(time.Hour))))
	arc.Issuer = user.ID.String()

	rrc := RESTClaims{
		IsAdmin: user.IsAdmin,
	}
	rrc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(a.c.JWT.RefreshLifespan * float32(time.Hour))))
	rrc.Issuer = user.ID.String()

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

	return &Tokens{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (a *restAuth) VerifyToken(token string, tokenType string) (*RESTClaims, error) {
	var secret string
	if tokenType == "access" {
		secret = a.c.JWT.AccessSecret
	} else if tokenType == "refresh" {
		secret = a.c.JWT.RefreshSecret
	} else {
		log.Panicln("invalid token type")
	}

	parsed, err := jwt.ParseWithClaims(token, &RESTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := parsed.Claims.(*RESTClaims)

	return claims, nil
}
