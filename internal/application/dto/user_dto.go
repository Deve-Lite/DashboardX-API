package dto

import (
	t "github.com/Deve-Lite/DashboardX-API-PoC/pkg/nullable"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type PasswordsUserRequest struct {
	Password    string `json:"password" binding:"required,min=6"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type GetUserResponse struct {
	ID       uuid.UUID `json:"id" binding:"required,uuid"`
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Theme    string    `json:"theme" binding:"required"`
	Language string    `json:"language" binding:"required"`
}

type UpdateUserRequest struct {
	Name     t.String `json:"name" binding:"nullmin=3" swaggertype:"string"`
	Email    t.String `json:"email" binding:"nullemail" swaggertype:"string"`
	Theme    t.String `json:"theme" swaggertype:"string"`
	Language t.String `json:"language" swaggertype:"string"`
}

type ChangeUserPasswordRequest struct {
	Password    string `json:"password" binding:"required,min=6" swaggertype:"string"`
	NewPassword string `json:"newPassword" binding:"required,min=6" swaggertype:"string"`
}

type DeleteUserRequest struct {
	Password string `json:"password" binding:"required,min=6" swaggertype:"string"`
}
