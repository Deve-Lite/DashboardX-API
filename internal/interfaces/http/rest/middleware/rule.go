package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Rule interface {
	LoggedIn(*gin.Context)
	ValidRefresh(ctx *gin.Context)
	ValidConfirm(ctx *gin.Context)
	ValidReset(ctx *gin.Context)
	ValidResetSubject(ctx *gin.Context)
}

type rule struct {
	a  application.RESTAuthService
	us application.UserService
}

func NewRule(a application.RESTAuthService, us application.UserService) Rule {
	return &rule{a, us}
}

func (r *rule) LoggedIn(ctx *gin.Context) {
	token := r.getToken(ctx)
	if token == "" {
		return
	}

	claims, err := r.a.VerifyToken(ctx, token, "access")
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

	claims, err := r.a.VerifyToken(ctx, token, "refresh")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
		return
	}

	if err := r.setUserContext(claims, ctx); err != nil {
		return
	}

	ctx.Next()
}

func (r *rule) ValidConfirm(ctx *gin.Context) {
	token := r.getToken(ctx)
	if token == "" {
		return
	}

	claims, err := r.a.VerifyConfirmToken(ctx, token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
		return
	}

	if err := r.setSubjectContext(claims, ctx); err != nil {
		return
	}

	ctx.Next()
}

func (r *rule) ValidReset(ctx *gin.Context) {
	token := r.getToken(ctx)
	if token == "" {
		return
	}

	claims, err := r.a.VerifyResetToken(ctx, token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
		return
	}

	if err := r.setSubjectContext(claims, ctx); err != nil {
		return
	}

	ctx.Next()
}

func (r *rule) ValidResetSubject(ctx *gin.Context) {
	subID, err := uuid.Parse(ctx.GetString("SubID"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	hashSubID, err := ctx.Cookie("rps")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
		return
	}

	if err := r.a.VerifyResetPasswordSubject(ctx, subID, hashSubID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(err))
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

	parts := strings.Split(bearer, " ")
	if len(parts) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(ae.ErrMissingAuthToken))
		return ""
	}

	if parts[1] == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ae.NewHTTPError(ae.ErrMissingAuthToken))
		return ""
	}

	return parts[1]
}

func (r *rule) setUserContext(claims *dto.RESTClaims, ctx *gin.Context) error {
	var userID uuid.UUID
	userID, err := uuid.Parse(claims.Subject)
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

	var JWTID uuid.UUID
	JWTID, err = uuid.Parse(claims.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return err
	}

	ctx.Set("UserID", user.ID.String())
	ctx.Set("IsAdmin", user.IsAdmin)
	ctx.Set("JWTID", JWTID.String())

	return nil
}

func (r *rule) setSubjectContext(claims *jwt.RegisteredClaims, ctx *gin.Context) error {
	var subID uuid.UUID
	subID, err := uuid.Parse(claims.Subject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return err
	}

	var JWTID uuid.UUID
	JWTID, err = uuid.Parse(claims.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return err
	}

	ctx.Set("SubID", subID.String())
	ctx.Set("JWTID", JWTID.String())

	return nil
}
