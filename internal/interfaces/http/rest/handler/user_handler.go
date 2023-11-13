package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API/internal/application/mapper"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	ae "github.com/Deve-Lite/DashboardX-API/pkg/errors"
	t "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Tokens(ctx *gin.Context)
	ConfirmAccount(ctx *gin.Context)
	ResendConfirmAccount(ctx *gin.Context)
	Login(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	ResetPasswordToken(ctx *gin.Context)
	ResetPasswordChange(ctx *gin.Context)
}

type userHandler struct {
	c  *config.Config
	us application.UserService
	m  mapper.UserMapper
}

func NewUserHandler(c *config.Config, us application.UserService, m mapper.UserMapper) UserHandler {
	return &userHandler{c, us, m}
}

// UserRegister godoc
//
//	@Summary		Register a new user
//	@Description	A link to confirm account will be sent to the provided email
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body	dto.CreateUserRequest	true	"Register input"
//	@Success		202
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		409	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Failure		503	{object}	errors.HTTPError
//	@Router			/users/register [post]
func (h *userHandler) Register(ctx *gin.Context) {
	body := &dto.CreateUserRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	_, err := h.us.PreCreate(ctx, h.m.CreateDTOToCreateModel(body))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ae.ErrEmailExists) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}

// UserConfirm godoc
//
//	@Summary		Confirm a newly registered account
//	@Description	Requires a valid confirm token sent in the Authorization header
//	@Security		BearerAuth
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		409	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/confirm-account [post]
func (h *userHandler) ConfirmAccount(ctx *gin.Context) {
	preUserID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	_, err = h.us.Create(ctx, preUserID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ae.ErrUserCreation) {
			code = http.StatusBadRequest
		} else if errors.Is(err, ae.ErrNoAwaitingConfirm) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

// UserResendConfirm godoc
//
//	@Summary		Send a new token to confirm an account
//	@Description	Sends a token to provided mailbox, if the account awaits to be confirmed
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body	dto.UserEmailRequest	true	"Resend confirm input"
//	@Success		202
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		409	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/confirm-account/resend [post]
func (h *userHandler) ResendConfirmAccount(ctx *gin.Context) {
	body := &dto.UserEmailRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	if err := h.us.SendConfirmToken(ctx, body.Email); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ae.ErrNoAwaitingConfirm) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}

// UserLogin godoc
//
//	@Summary	Login a user
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		data	body		dto.LoginUserRequest	true	"Login data"
//	@Success	200		{object}	dto.Tokens
//	@Failure	400		{object}	errors.HTTPError
//	@Failure	404		{object}	errors.HTTPError
//	@Failure	409		{object}	errors.HTTPError
//	@Failure	500		{object}	errors.HTTPError
//	@Router		/users/login [post]
func (h *userHandler) Login(ctx *gin.Context) {
	body := &dto.LoginUserRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	tokens, err := h.us.Login(ctx, h.m.LoginDTOToModel(body))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ae.ErrUserNotFound) {
			code = http.StatusNotFound
		} else if errors.Is(err, ae.ErrInvalidPassword) {
			code = http.StatusBadRequest
		} else if errors.Is(err, ae.ErrConfirmationRequired) {
			code = http.StatusConflict
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.m.TokenModelToTokenDTO(tokens))
}

// UserMeTokens godoc
//
//	@Summary		Generate a new pair of user tokens
//	@Description	Requires a valid refresh token sent in the Authorization header
//	@Security		BearerAuth
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.Tokens
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/me/tokens [post]
func (h *userHandler) Tokens(ctx *gin.Context) {
	userID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	tokens, err := h.us.GetTokens(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.m.TokenModelToTokenDTO(tokens))
}

// UserRefresh godoc
//
//	@Deprecated
//	@Summary		Generate a new pair of user tokens
//	@Description	Requires a valid refresh token sent in the Authorization header
//	@Security		BearerAuth
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.Tokens
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/refresh [post]
func (h *userHandler) Refresh(ctx *gin.Context) {}

// UserGetMe godoc
//
//	@Summary	Get a logged in user
//	@Security	BearerAuth
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	dto.GetUserResponse
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/users/me [get]
func (h *userHandler) Get(ctx *gin.Context) {
	userID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	var user *domain.User
	user, err = h.us.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, ae.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.m.ModelToDTO(user))
}

// UserUpdateMe godoc
//
//	@Summary	Update a logged in user
//	@Security	BearerAuth
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		data	body	dto.UpdateUserRequest	true	"Update data"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	409	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/users/me [patch]
func (h *userHandler) Update(ctx *gin.Context) {
	userID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	body := &dto.UpdateUserRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	user := h.m.UpdateDTOToUpdateModel(body)
	user.ID = userID

	err = h.us.Update(ctx, user)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ae.ErrEmailExists) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UserDeleteMe godoc
//
//	@Summary	Delete a logged in user
//	@Security	BearerAuth
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		data	body	dto.DeleteUserRequest	true	"Delete data"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/users/me [delete]
func (h *userHandler) Delete(ctx *gin.Context) {
	userID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	body := &dto.DeleteUserRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	err = h.us.Verify(ctx, userID, body.Password)
	if err != nil {
		if errors.Is(err, ae.ErrInvalidPassword) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	err = h.us.Delete(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UserChangePassword godoc
//
//	@Summary	Change a password of logged in user
//	@Security	BearerAuth
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		data	body	dto.ChangeUserPasswordRequest	true	"Change passoword data"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/users/me/password [patch]
func (h *userHandler) ChangePassword(ctx *gin.Context) {
	userID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	body := &dto.ChangeUserPasswordRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	err = h.us.Verify(ctx, userID, body.Password)
	if err != nil {
		if errors.Is(err, ae.ErrInvalidPassword) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	user := &domain.UpdateUser{
		ID:       userID,
		Password: t.NewString(body.NewPassword, false, true),
	}

	err = h.us.Update(ctx, user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UserResetPasswordToken godoc
//
//	@Summary		Call an action to reset user's password
//	@Description	User receives a token at a given email
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body	dto.UserEmailRequest	true	"Reset passoword data"
//	@Success		202
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/reset-password [post]
func (h *userHandler) ResetPasswordToken(ctx *gin.Context) {
	body := &dto.UserEmailRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	age := int(time.Duration(h.c.JWT.ResetLifespanMinutes * float32(time.Minute)).Milliseconds())

	hashSubID, err := h.us.SendResetToken(ctx, body.Email)
	if err != nil {
		// return successful response when user does not exist to do not let scan the API
		if !errors.Is(err, ae.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
			return
		}
	}

	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("rps", hashSubID, age, "/api/v1/users/reset-password", h.c.Server.Domain, true, true)

	ctx.Status(http.StatusAccepted)
}

// UserResetPasswordChange godoc
//
//	@Summary		Set a new password for a user
//	@Description	Requires a token sent in the Authorization header to verify password's change action
//	@Security		BearerAuth
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body	dto.ResetUserPasswordRequest	true	"Reset passoword data"
//	@Success		204
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/reset-password [patch]
func (h *userHandler) ResetPasswordChange(ctx *gin.Context) {
	defer func() {
		ctx.SetSameSite(http.SameSiteStrictMode)
		ctx.SetCookie("rps", "", -1, "/", h.c.Server.Domain, true, true)
	}()
	subID, err := h.resolveSubject(ctx)
	if err != nil {
		return
	}

	body := &dto.ResetUserPasswordRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	if err := h.us.ResetPassword(ctx, subID, body.Password); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(ae.ErrUnexpected))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *userHandler) resolveSubject(ctx *gin.Context) (uuid.UUID, error) {
	ID := uuid.Nil
	for _, key := range []string{"UserID", "SubID"} {
		v, _ := ctx.Get(key)
		if v == nil {
			continue
		}

		var err error
		ID, err = uuid.Parse(v.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
			return uuid.Nil, err
		}
	}

	if ID == uuid.Nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(ae.ErrUnexpected))
		return uuid.Nil, ae.ErrUnexpected
	}

	return ID, nil
}
