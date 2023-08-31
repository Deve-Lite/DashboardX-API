package handler

import (
	"errors"
	"net/http"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	_ "github.com/Deve-Lite/DashboardX-API-PoC/internal/interfaces/http/rest/auth"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/mapper"
	ae "github.com/Deve-Lite/DashboardX-API-PoC/pkg/errors"
	t "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	Login(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type userHandler struct {
	us application.UserService
	m  mapper.UserMapper
}

func NewUserHandler(us application.UserService, m mapper.UserMapper) UserHandler {
	return &userHandler{us, m}
}

// UserRegister godoc
//
//	@Summary	Register a new user
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		data	body	dto.CreateUserRequest	true	"Register input"
//	@Success	201
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	409	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/users/register [post]
func (h *userHandler) Register(ctx *gin.Context) {
	body := &dto.CreateUserRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	_, err := h.us.Create(ctx, h.m.CreateDTOToCreateModel(body))
	if err != nil {
		var code int
		if errors.Is(err, ae.ErrEmailExists) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

// UserLogin godoc
//
//	@Summary	Login a user
//	@Tags		Users
//	@Accept		json
//	@Produce	json
//	@Param		data	body		dto.LoginUserRequest	true	"Login data"
//	@Success	200		{object}	auth.Tokens
//	@Failure	400		{object}	errors.HTTPError
//	@Failure	404		{object}	errors.HTTPError
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
		var code int
		if errors.Is(err, ae.ErrUserNotFound) {
			code = http.StatusNotFound
		} else if errors.Is(err, ae.ErrInvalidPassword) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.m.TokenModelToTokenDTO(tokens))
}

// UserRefresh godoc
//
//	@Summary		Generate a new pair of user tokens
//	@Description	Requires a valid refresh token sent in the Authorization header
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	auth.Tokens
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/users/refresh [post]
func (h *userHandler) Refresh(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return
	}

	tokens, err := h.us.Refresh(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.m.TokenModelToTokenDTO(tokens))
}

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
	userID, err := h.getUserID(ctx)
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
	userID, err := h.getUserID(ctx)
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
		var code int
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
	userID, err := h.getUserID(ctx)
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
	userID, err := h.getUserID(ctx)
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

func (h *userHandler) getUserID(ctx *gin.Context) (uuid.UUID, error) {
	userID, err := uuid.Parse(ctx.MustGet("UserID").(string))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return uuid.Nil, err
	}

	return userID, nil
}
