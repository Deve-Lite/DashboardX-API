package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/auth"
	ae "github.com/Deve-Lite/DashboardX-API-PoC/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Rule interface {
	LoggedIn(*gin.Context)
	ValidRefresh(ctx *gin.Context)
}

type rule struct {
	a  auth.RESTAuth
	us application.UserService
}

func NewRule(a auth.RESTAuth, us application.UserService) Rule {
	return &rule{a, us}
}

func (r *rule) LoggedIn(ctx *gin.Context) {
	token := r.getToken(ctx)
	if token == "" {
		return
	}

	claims, err := r.a.VerifyToken(token, "access")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
		return
	}

	if err := r.setUserContext(claims, ctx); err != nil {
		return
	}

	ctx.Next()
}

func (r *rule) ValidRefresh(ctx *gin.Context) {
	token := r.getToken(ctx)
	if token == "" {
		return
	}

	claims, err := r.a.VerifyToken(token, "refresh")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
		return
	}

	if err := r.setUserContext(claims, ctx); err != nil {
		return
	}

	ctx.Next()
}

func (r *rule) getToken(ctx *gin.Context) string {
	bearer := ctx.GetHeader("authorization")
	if bearer == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(ae.ErrMissingAuthToken))
		return ""
	}

	token := strings.Split(bearer, " ")[1]
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(ae.ErrMissingAuthToken))
		return ""
	}

	return token
}

func (r *rule) setUserContext(claims *auth.RESTClaims, ctx *gin.Context) error {
	var userID uuid.UUID
	userID, err := uuid.Parse(claims.Issuer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return err
	}

	var user *domain.User
	user, err = r.us.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, ae.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return err
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return err
	}

	ctx.Set("UserID", user.ID.String())
	ctx.Set("IsAdmin", user.IsAdmin)

	return nil
}
